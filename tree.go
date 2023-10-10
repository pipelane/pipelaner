package pipelane

import (
	"context"
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
)

var (
	ErrInputsNotFound       = errors.New("ErrInputsNotFound")
	ErrLaneNameMustBeUnique = errors.New("ErrLaneNameMustBeUnique")
	ErrInvalidConfig        = errors.New("ErrInvalidConfig")
)

func ErrLaneWithoutSink(s string) error {
	return errors.New(fmt.Sprintf("ErrLaneWithoutSink: %s", s))
}

type Transformer interface {
	Transform(ctx context.Context, val any) any
}
type Sink interface {
	Sink(ctx context.Context, val any)
}
type Generator interface {
	Generate(ctx context.Context) any
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
		return i.Cfg.LaneType == lt
	})
}

func (t *TreeLanes) mapWithInputs() map[string][]*LaneItem {
	inputs := map[string][]*LaneItem{}
	for _, val := range t.Items {
		if val.Cfg.Input == nil {
			continue
		}
		var arr []*LaneItem
		if i, ok := inputs[*val.Cfg.Input]; ok {
			arr = i
		}
		arr = append(arr, val)
		inputs[*val.Cfg.Input] = arr
	}
	return inputs
}

func (t *TreeLanes) append(val *LaneItem) {
	t.Items[val.Cfg.Name] = val
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
	var c map[string]any
	_, err := toml.DecodeFile(file, &c)
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

	err = flat(ctx, InputType, cfg.Input, lanes)
	if err != nil {
		return nil, err
	}

	err = flat(ctx, LaneType, cfg.Lane, lanes)
	if err != nil {
		return nil, err
	}

	err = flat(ctx, SinkType, cfg.Sink, lanes)
	if err != nil {
		return nil, err
	}

	if len(cfg.Sink)+len(cfg.Lane)+len(cfg.Input) != len(lanes.Items) {
		return nil, ErrLaneNameMustBeUnique
	}
	lanes.connect()
	if err = lanes.run(dataSource); err != nil {
		return nil, err
	}
	return lanes, nil
}

func flat(
	ctx context.Context,
	types LaneTypes,
	input map[string]any,
	output *TreeLanes,
) error {
	for k, v := range input {
		val, ok := v.(map[string]any)
		if !ok {
			return ErrInvalidConfig
		}
		cfg, err := NewBaseConfigWithTypeAndExtended(types, val)
		if err != nil {
			return err
		}

		cfg.Name = k
		p := NewLaneItem(
			ctx,
			cfg,
		)
		output.append(p)
	}
	return nil
}

func (t *TreeLanes) run(dataSource DataSource) error {
	inputs := t.filterByType(InputType)
	if err := t.validateOutputs(inputs); err != nil {
		return err
	}
	transforms := t.filterByType(LaneType)
	if err := t.validateOutputs(transforms); err != nil {
		return err
	}
	sinks := t.filterByType(SinkType)
	for i := range transforms {
		item := transforms[i]
		tr := dataSource.Transforms[item.Cfg.SourceName]
		item.subscriber.Transform(tr.Transform)
		go item.subscriber.run()
	}
	for i := range sinks {
		item := sinks[i]
		si := dataSource.Sinks[item.Cfg.SourceName]
		item.subscriber.Sink(si.Sink)
		go item.subscriber.run()
	}
	for i := range inputs {
		item := inputs[i]
		generator := dataSource.Generators[item.Cfg.SourceName]
		item.Gen(generator.Generate)
		go item.subscriber.run()
		go item.subscriber.Generate()
	}
	return nil
}

func (t *TreeLanes) connect() {
	allWithInputs := t.mapWithInputs()
	for i := range t.Items {
		input := t.Items[i]
		output, ok := allWithInputs[input.Cfg.Name]
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
			return ErrLaneWithoutSink(l.Cfg.Name)
		}
	}
	return nil
}
