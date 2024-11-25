package pipelaner

import (
	"context"
	"fmt"

	"github.com/pipelane/pipelaner/gen/components"
	"github.com/pipelane/pipelaner/internal/pipeline"
)

type Pipelaner struct {
	Pipelines []*pipeline.Pipeline
}

func NewPipelaner(configs []*components.Pipeline) (*Pipelaner, error) {
	pipelines := make([]*pipeline.Pipeline, 0, len(configs))

	for _, cfg := range configs {
		p, err := pipeline.NewPipeline(cfg)
		if err != nil {
			return nil, fmt.Errorf("create pipeline '%s': %w", cfg.Name, err)
		}
		pipelines = append(pipelines, p)
	}

	return &Pipelaner{
		Pipelines: pipelines,
	}, nil
}

func (p *Pipelaner) Run(ctx context.Context) error {
	for _, pipe := range p.Pipelines {
		if err := pipe.Run(ctx); err != nil {
			return err
		}
	}
	return nil
}
