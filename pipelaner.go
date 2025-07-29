/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"fmt"

	"github.com/pipelane/pipelaner/gen/components"
	logCfg "github.com/pipelane/pipelaner/gen/settings/logger"
	"github.com/pipelane/pipelaner/internal/logger"
	pipelines "github.com/pipelane/pipelaner/pipeline"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type pipeline interface {
	Run(ctx context.Context) error
}

type Pipelaner struct {
	Pipelines []pipeline
	l         *zerolog.Logger
}

func NewPipelaner(
	configs []components.Pipeline,
	loggerCfg logCfg.Config,
	metricsEnabled, gcAfterProcess bool,
) (*Pipelaner, error) {
	pl := make([]pipeline, 0, len(configs))
	l, err := logger.NewLoggerWithCfg(loggerCfg)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	for _, cfg := range configs {
		p, e := pipelines.NewPipeline(cfg, l, metricsEnabled, gcAfterProcess)
		if e != nil {
			return nil, fmt.Errorf("create pipeline '%s': %w", cfg.Name, e)
		}
		pl = append(pl, p)
	}

	return &Pipelaner{
		Pipelines: pl,
		l:         l,
	}, nil
}

func (p *Pipelaner) Run(ctx context.Context) error {
	gr := errgroup.Group{}
	for _, pipe := range p.Pipelines {
		gr.Go(func() error {
			err := pipe.Run(ctx)
			if err != nil {
				p.l.Panic().Err(err).Msg("pipeline run")
				return err
			}
			return nil
		})
	}
	return gr.Wait()
}
