/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrInputsNotFound       = errors.New("ErrInputsNotFound")
	ErrLaneNameMustBeUnique = errors.New("ErrLaneNameMustBeUnique")
	ErrInvalidConfig        = errors.New("ErrInvalidConfig")
)

func ErrLaneWithoutSink(s string) error {
	return fmt.Errorf("ErrLaneWithoutSink: %s", s)
}

type Init interface {
	Init(ctx *Context) error
}

type Map interface {
	Init
	New() Map
	Map(ctx *Context, val any) any
}
type Sink interface {
	Init
	Sink(ctx *Context, val any)
}
type Generator interface {
	Init
	Generate(ctx *Context, input chan<- any)
}

type TreeLanes struct {
	Items map[string]*LaneItem
}

func (t *TreeLanes) filter(f func(i *LaneItem) bool) []*LaneItem {
	var res []*LaneItem
	for _, v := range t.Items {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}

func (t *TreeLanes) filterByType(lt LaneTypes) []*LaneItem {
	return t.filter(func(i *LaneItem) bool {
		return i.cfg.LaneType == lt
	})
}

func (t *TreeLanes) mapWithInputs() map[string][]*LaneItem {
	inputs := map[string][]*LaneItem{}
	for _, val := range t.Items {
		if len(val.cfg.Inputs) == 0 {
			continue
		}
		for i := range val.cfg.Inputs {
			input := val.cfg.Inputs[i]
			var arr []*LaneItem
			if v, ok := inputs[input]; ok {
				arr = v
			}
			arr = append(arr, val)
			inputs[input] = arr
		}
	}
	return inputs
}

func (t *TreeLanes) append(val *LaneItem) {
	t.Items[val.cfg.Name] = val
}

func newPipelinesTree() *TreeLanes {
	return &TreeLanes{
		Items: map[string]*LaneItem{},
	}
}

func NewTreeFrom(
	ctx context.Context,
	dataSource DataSource,
	file string,
) (*TreeLanes, error) {
	c, err := readToml(file)
	if err != nil {
		return nil, err
	}
	a, err := newPipelinesTreeMapWith(ctx, dataSource, c)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func newPipelinesTreeMapWith(
	ctx context.Context,
	dataSource DataSource,
	c map[string]any,
) (*TreeLanes, error) {
	lanes := newPipelinesTree()
	cfg, err := newConfig(c)
	if err != nil {
		return nil, err
	}
	if len(cfg.Input) == 0 {
		return nil, ErrInputsNotFound
	}

	err = flat(InputType, cfg.Input, lanes)
	if err != nil {
		return nil, err
	}

	err = flat(MapType, cfg.Map, lanes)
	if err != nil {
		return nil, err
	}

	err = flat(SinkType, cfg.Sink, lanes)
	if err != nil {
		return nil, err
	}

	if len(cfg.Sink)+len(cfg.Map)+len(cfg.Input) != len(lanes.Items) {
		return nil, ErrLaneNameMustBeUnique
	}
	lanes.connect(ctx)
	if err = lanes.run(ctx, dataSource); err != nil {
		return nil, err
	}
	return lanes, nil
}

func flat(
	types LaneTypes,
	input map[string]any,
	output *TreeLanes,
) error {
	for k, v := range input {
		val, ok := v.(map[string]any)
		if !ok {
			return ErrInvalidConfig
		}
		cfg, err := NewBaseConfigWithTypeAndExtended(types, k, val)
		if err != nil {
			return err
		}
		p := NewLaneItem(
			cfg,
		)
		output.append(p)
	}
	return nil
}

func (t *TreeLanes) run(ctx context.Context, dataSource DataSource) error {
	inputs := t.filterByType(InputType)
	if err := t.validateOutputs(inputs); err != nil {
		return err
	}
	transforms := t.filterByType(MapType)
	if err := t.validateOutputs(transforms); err != nil {
		return err
	}
	sinks := t.filterByType(SinkType)
	for i := range transforms {
		item := transforms[i]
		c := &Context{
			ctx:      ctx,
			laneItem: item,
		}
		tr := dataSource.Maps[c.LaneItem().Config().SourceName]
		t := tr.New()
		err := t.Init(c)
		if err != nil {
			return err
		}
		item.runLoop.setMap(t.Map)
		item.runLoop.run()
	}
	for i := range sinks {
		item := sinks[i]
		c := &Context{
			ctx:      ctx,
			laneItem: item,
		}
		si := dataSource.Sinks[c.LaneItem().Config().SourceName]
		err := si.Init(c)
		if err != nil {
			return err
		}
		item.runLoop.setSink(si.Sink)
		item.runLoop.run()
	}
	for i := range inputs {
		item := inputs[i]
		c := &Context{
			ctx:      ctx,
			laneItem: item,
		}
		generator := dataSource.Generators[c.LaneItem().Config().SourceName]
		err := generator.Init(c)
		if err != nil {
			return err
		}
		item.runLoop.setGenerator(generator.Generate)
		item.runLoop.run()
		item.runLoop.receive()
	}
	return nil
}

func (t *TreeLanes) connect(ctx context.Context) {
	allWithInputs := t.mapWithInputs()
	inputs := t.filterByType(InputType)
	for i := range inputs {
		input := inputs[i]
		input.runLoop.setContext(NewContext(ctx, input))
	}
	for i := range t.Items {
		input := t.Items[i]
		output, ok := allWithInputs[input.cfg.Name]
		if !ok {
			continue
		}
		for j := range output {
			out := output[j]
			input.Subscribe(out)
		}
	}
}

func (t *TreeLanes) validateOutputs(lanes []*LaneItem) error {
	for i := range lanes {
		l := lanes[i]
		if len(l.outputPipelines) == 0 {
			return ErrLaneWithoutSink(l.cfg.Name)
		}
	}
	return nil
}
