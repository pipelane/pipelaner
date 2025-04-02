package source

import (
	"fmt"

	"github.com/LastPossum/kamino"
	"github.com/pipelane/pipelaner/pipeline/components"
)

type Inputs map[string]components.Input
type Transforms map[string]components.Transform
type Sequencers map[string]components.Sequencer
type Sinks map[string]components.Sink

type dataSources struct {
	Inputs     Inputs
	Transforms Transforms
	Sinks      Sinks
	Sequencers Sequencers
}

func newDataSources() *dataSources {
	return &dataSources{
		Inputs:     make(Inputs),
		Transforms: make(Transforms),
		Sinks:      make(Sinks),
		Sequencers: make(Sequencers),
	}
}

var dataSource = newDataSources()

func RegisterSink(name string, sink components.Sink) {
	dataSource.Sinks[name] = sink
}

func GetSink(name string) (components.Sink, error) {
	sink, ok := dataSource.Sinks[name]
	if !ok {
		return nil, fmt.Errorf("sink: %s is not registered", name)
	}
	s, err := kamino.Clone(sink)
	if err != nil {
		return nil, fmt.Errorf("clone sink: %s", name)
	}
	return s, nil
}

func RegisterInput(name string, input components.Input) {
	dataSource.Inputs[name] = input
}

func GetInput(name string) (components.Input, error) {
	input, ok := dataSource.Inputs[name]
	if !ok {
		return nil, fmt.Errorf("input: %s is not registered", name)
	}
	i, err := kamino.Clone(input)
	if err != nil {
		return nil, fmt.Errorf("clone input: %s", name)
	}
	return i, nil
}

func RegisterTransform(name string, transform components.Transform) {
	dataSource.Transforms[name] = transform
}

func GetTransform(name string) (components.Transform, error) {
	transform, ok := dataSource.Transforms[name]
	if !ok {
		return nil, fmt.Errorf("transform: %s is not registered", name)
	}
	t, err := kamino.Clone(transform)
	if err != nil {
		return nil, fmt.Errorf("clone transform: %s", name)
	}
	return t, nil
}

func RegisterSequencer(name string, sequencer components.Sequencer) {
	dataSource.Sequencers[name] = sequencer
}

func GetSequencer(name string) (components.Sequencer, error) {
	sequencer, ok := dataSource.Sequencers[name]
	if !ok {
		return nil, fmt.Errorf("sequencer: %s is not registered", name)
	}
	t, err := kamino.Clone(sequencer)
	if err != nil {
		return nil, fmt.Errorf("clone sequencer: %s", name)
	}
	return t, nil
}
