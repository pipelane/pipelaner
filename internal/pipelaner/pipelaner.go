package pipelaner

import (
	"context"
	"fmt"

	"github.com/pipelane/pipelaner/gen/components"
	logCfg "github.com/pipelane/pipelaner/gen/settings/logger"
	"github.com/pipelane/pipelaner/internal/logger"
	pipelineimpl "github.com/pipelane/pipelaner/internal/pipeline"
)

type pipeline interface {
	Run(ctx context.Context) error
}

type Pipelaner struct {
	Pipelines []pipeline
}

func NewPipelaner(
	configs []*components.Pipeline,
	loggerCfg *logCfg.LoggerConfig,
	metricsEnabled, gcAfterProcess bool,
) (*Pipelaner, error) {
	pipelines := make([]pipeline, 0, len(configs))
	l, err := logger.NewLoggerWithCfg(loggerCfg)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	for _, cfg := range configs {
		p, err := pipelineimpl.NewPipeline(cfg, l, metricsEnabled, gcAfterProcess)
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
