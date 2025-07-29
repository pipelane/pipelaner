/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package pipeline

import (
	"context"
	"fmt"
	"slices"

	config "github.com/pipelane/pipelaner/gen/components"
	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/pipeline/node"
	"github.com/rs/zerolog"
)

type (
	inputNode interface {
		Run(ctx context.Context) error
		AddOutputChannel(ch chan any)

		GetName() string
		GetOutputBufferSize() uint
	}

	transformNode interface {
		Run() error
		AddInputChannel(ch chan any)
		AddOutputChannel(ch chan any)

		GetName() string
		GetInputs() []string
		GetOutputBufferSize() uint
	}

	sinkNode interface {
		Run() error
		AddInputChannel(ch chan any)

		GetInputs() []string
	}
)

type Pipeline struct {
	name       string
	inputs     []inputNode
	transforms []transformNode
	sinks      []sinkNode
}

func NewPipeline(
	cfg config.Pipeline,
	logger *zerolog.Logger,
	enableMetrics, startGCAfterProcess bool,
) (*Pipeline, error) {
	p := &Pipeline{
		name: cfg.Name,
	}

	var opts []node.Option
	if enableMetrics {
		opts = append(opts, node.WithMetrics())
	}
	if startGCAfterProcess {
		opts = append(opts, node.WithCallGC())
	}

	if err := p.initNodes(cfg, logger, opts...); err != nil {
		return nil, err
	}
	p.connectNodes()
	return p, nil
}

func (p *Pipeline) Run(ctx context.Context) error {
	for _, sinkNode := range p.sinks {
		if err := sinkNode.Run(); err != nil {
			return fmt.Errorf("run sink node: %w", err)
		}
	}

	for _, transformNode := range p.transforms {
		if err := transformNode.Run(); err != nil {
			return fmt.Errorf("run transform node: %w", err)
		}
	}

	for _, inputNode := range p.inputs {
		if err := inputNode.Run(ctx); err != nil {
			return fmt.Errorf("run input node: %w", err)
		}
	}
	<-ctx.Done()
	return nil
}

func (p *Pipeline) initNodes(cfg config.Pipeline, logger *zerolog.Logger, opts ...node.Option) error {
	if err := p.initSinkNodes(cfg.Sinks, logger, opts...); err != nil {
		return fmt.Errorf("init sink nodes: %w", err)
	}
	if err := p.initTransformNodes(cfg.Transforms, logger, opts...); err != nil {
		return fmt.Errorf("init transform nodes: %w", err)
	}
	if err := p.initInputNodes(cfg.Inputs, logger, opts...); err != nil {
		return fmt.Errorf("init input nodes: %w", err)
	}
	return nil
}

func (p *Pipeline) connectNodes() {
	p.connectSinkNodes()
	p.connectTransformNodes()
}

func (p *Pipeline) initSinkNodes(
	sinkConfigs []sink.Sink,
	logger *zerolog.Logger,
	opts ...node.Option,
) error {
	sinkNodes := make([]sinkNode, 0, len(sinkConfigs))
	for _, cfg := range sinkConfigs {
		sinkNode, err := node.NewSink(cfg, logger, opts...)
		if err != nil {
			return fmt.Errorf("new sink node: %w", err)
		}
		sinkNodes = append(sinkNodes, sinkNode)
	}
	p.sinks = sinkNodes
	return nil
}

func (p *Pipeline) initTransformNodes(
	transformConfigs []transform.Transform,
	logger *zerolog.Logger,
	opts ...node.Option,
) error {
	transformNodes := make([]transformNode, 0, len(transformConfigs))
	var transformsNode transformNode
	var err error
	for _, cfg := range transformConfigs {
		switch cfg.(type) {
		case transform.Sequencer:
			transformsNode, err = node.NewSequencer(cfg, logger, opts...)
		case transform.Transform:
			transformsNode, err = node.NewTransform(cfg, logger, opts...)
		}
		if err != nil {
			return fmt.Errorf("new transform node: %w", err)
		}
		transformNodes = append(transformNodes, transformsNode)
	}
	p.transforms = transformNodes
	return nil
}

func (p *Pipeline) initInputNodes(
	inputs []input.Input,
	logger *zerolog.Logger,
	opts ...node.Option,
) error {
	inputNodes := make([]inputNode, 0, len(inputs))
	for _, cfg := range inputs {
		inputNode, err := node.NewInput(cfg, logger, opts...)
		if err != nil {
			return fmt.Errorf("new input node: %w", err)
		}
		inputNodes = append(inputNodes, inputNode)
	}
	p.inputs = inputNodes
	return nil
}

func (p *Pipeline) connectSinkNodes() {
	for _, sinkNode := range p.sinks {
		for _, transformNode := range p.transforms {
			if slices.Contains(sinkNode.GetInputs(), transformNode.GetName()) {
				ch := make(chan any, transformNode.GetOutputBufferSize())
				sinkNode.AddInputChannel(ch)
				transformNode.AddOutputChannel(ch)
			}
		}
		for _, inputNode := range p.inputs {
			if slices.Contains(sinkNode.GetInputs(), inputNode.GetName()) {
				ch := make(chan any, inputNode.GetOutputBufferSize())
				sinkNode.AddInputChannel(ch)
				inputNode.AddOutputChannel(ch)
			}
		}
	}
}

func (p *Pipeline) connectTransformNodes() {
	for _, tNode := range p.transforms {
		for _, transformNode := range p.transforms {
			if slices.Contains(tNode.GetInputs(), transformNode.GetName()) {
				ch := make(chan any, transformNode.GetOutputBufferSize())
				tNode.AddInputChannel(ch)
				transformNode.AddOutputChannel(ch)
			}
		}

		for _, inputNode := range p.inputs {
			if slices.Contains(tNode.GetInputs(), inputNode.GetName()) {
				ch := make(chan any, inputNode.GetOutputBufferSize())
				tNode.AddInputChannel(ch)
				inputNode.AddOutputChannel(ch)
			}
		}
	}
}
