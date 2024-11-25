package pipeline

import (
	"context"
	"fmt"
	"slices"

	config "github.com/pipelane/pipelaner/gen/components"
	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/internal/pipeline/node"
)

type (
	inputNode interface {
		Run(ctx context.Context) error
		AddOutputChannel(ch chan any)

		GetName() string
		GetOutputBufferSize() int
	}

	transformNode interface {
		Run() error
		AddInputChannel(ch chan any)
		AddOutputChannel(ch chan any)

		GetName() string
		GetInputs() []string
		GetOutputBufferSize() int
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

func NewPipeline(cfg *config.Pipeline) (*Pipeline, error) {
	p := &Pipeline{
		name: cfg.Name,
	}

	if err := p.initNodes(cfg); err != nil {
		return nil, err
	}

	if err := p.connectNodes(); err != nil {
		return nil, err
	}

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

func (p *Pipeline) initNodes(cfg *config.Pipeline) error {
	if err := p.initSinkNodes(cfg.Sinks); err != nil {
		return fmt.Errorf("init sink nodes: %w", err)
	}
	if err := p.initTransformNodes(cfg.Maps); err != nil {
		return fmt.Errorf("init transform nodes: %w", err)
	}
	if err := p.initInputNodes(cfg.Inputs); err != nil {
		return fmt.Errorf("init input nodes: %w", err)
	}
	return nil
}

func (p *Pipeline) connectNodes() error {
	if err := p.connectSinkNodes(); err != nil {
		return fmt.Errorf("connect sink nodes: %w", err)
	}
	if err := p.connectTransformNodes(); err != nil {
		return fmt.Errorf("connect transform nodes: %w", err)
	}
	return nil
}

func (p *Pipeline) initSinkNodes(sinkConfigs []sink.Sink) error {
	sinkNodes := make([]sinkNode, 0, len(sinkConfigs))
	for _, cfg := range sinkConfigs {
		// todo add logger and options
		sinkNode, err := node.NewSink(cfg, nil)
		if err != nil {
			return fmt.Errorf("new sink node: %w", err)
		}
		sinkNodes = append(sinkNodes, sinkNode)
	}
	p.sinks = sinkNodes
	return nil
}

func (p *Pipeline) initTransformNodes(transformConfigs []transform.Transform) error {
	transformNodes := make([]transformNode, 0, len(transformConfigs))
	for _, cfg := range transformConfigs {
		// todo add logger and options
		transformNode, err := node.NewTransform(cfg, nil)
		if err != nil {
			return fmt.Errorf("new transform node: %w", err)
		}
		transformNodes = append(transformNodes, transformNode)
	}
	p.transforms = transformNodes
	return nil
}

func (p *Pipeline) initInputNodes(inputs []input.Input) error {
	inputNodes := make([]inputNode, 0, len(inputs))
	for _, cfg := range inputs {
		inputNode, err := node.NewInput(cfg, nil)
		if err != nil {
			return fmt.Errorf("new input node: %w", err)
		}
		inputNodes = append(inputNodes, inputNode)
	}
	p.inputs = inputNodes
	return nil
}

func (p *Pipeline) connectSinkNodes() error {
	for _, sinkNode := range p.sinks {
		// sink нода может иметь связь с transform и input нодами
		// 1. проверяем transform ноды на связь с текущим sink
		for _, transformNode := range p.transforms {
			if slices.Contains(sinkNode.GetInputs(), transformNode.GetName()) {
				ch := make(chan any, transformNode.GetOutputBufferSize())
				sinkNode.AddInputChannel(ch)
				transformNode.AddOutputChannel(ch)
			}
		}
		// 2. проверяем input ноды на связь с текущим input
		for _, inputNode := range p.inputs {
			if slices.Contains(sinkNode.GetInputs(), inputNode.GetName()) {
				ch := make(chan any, inputNode.GetOutputBufferSize())
				sinkNode.AddInputChannel(ch)
				inputNode.AddOutputChannel(ch)
			}
		}
	}
	return nil
}

func (p *Pipeline) connectTransformNodes() error {
	for _, tNode := range p.transforms {
		// transform нода может быть связана с transform и input нодами
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
	return nil
}
