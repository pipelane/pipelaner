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
	ErrUnknownItem          = errors.New("ErrUnknownItem")
	ErrLaneNameMustBeUnique = errors.New("ErrLaneNameMustBeUnique")
	ErrInvalidConfig        = errors.New("ErrInvalidConfig")
	ErrUnknownGenerator     = errors.New("ErrUnknownGenerator")
	ErrUnknownMap           = errors.New("ErrUnknownMap")
	ErrUnknownSink          = errors.New("ErrUnknownSink")
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

func (t *TreeLanes) sortedByType() []*LaneItem {
	l := t.filterByType(InputType)
	l = append(l, t.filterByType(MapType)...)
	l = append(l, t.filterByType(SinkType)...)
	return l
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
	file string,
) (*TreeLanes, error) {
	c, err := ReadToml(file)
	if err != nil {
		return nil, err
	}
	a, err := newPipelinesTreeMapWith(ctx, c)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func newPipelinesTreeMapWith(
	ctx context.Context,
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

	// set logger into context
	logger, err := initLogger(cfg)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, logKey, logger)

	if len(cfg.Sink)+len(cfg.Map)+len(cfg.Input) != len(lanes.Items) {
		return nil, ErrLaneNameMustBeUnique
	}
	lanes.makeTree(ctx)
	if err = lanes.run(ctx); err != nil {
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

func (t *TreeLanes) run(ctx context.Context) error {
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
		tr, ok := dataSource.Maps[c.LaneItem().Config().SourceName]
		if !ok {
			return ErrUnknownMap
		}
		nT := tr.New()
		err := nT.Init(c)
		if err != nil {
			return err
		}
		item.runLoop.setMap(nT.Map)
		item.runLoop.run()
	}
	for i := range sinks {
		item := sinks[i]
		c := &Context{
			ctx:      ctx,
			laneItem: item,
		}
		si, ok := dataSource.Sinks[c.LaneItem().Config().SourceName]
		if !ok {
			return ErrUnknownSink
		}
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
		generator, ok := dataSource.Generators[c.LaneItem().Config().SourceName]
		if !ok {
			return ErrUnknownGenerator
		}
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

func (t *TreeLanes) makeTree(ctx context.Context) {
	allWithInputs := t.mapWithInputs()
	inputs := t.filterByType(InputType)
	for i := range inputs {
		input := inputs[i]
		input.runLoop.setContext(NewContext(ctx, input))
	}
	sorted := t.sortedByType()
	for i := range sorted {
		input := sorted[i]
		output, ok := allWithInputs[input.cfg.Name]
		if !ok {
			continue
		}
		for j := range output {
			out := output[j]
			input.addOutputs(out)
		}
	}
	t.subscribeRecursive(inputs)
}
func (t *TreeLanes) subscribeRecursive(inputs []*LaneItem) {
	if len(inputs) == 0 {
		return
	}
	for _, input := range inputs {
		for j := range input.outputPipelines {
			output := input.outputPipelines[j]
			input.Subscribe(output)
		}
		t.subscribeRecursive(input.outputPipelines)
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
