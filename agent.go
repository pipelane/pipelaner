/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"fmt"

	config "github.com/pipelane/pipelaner/gen/pipelaner"
	"github.com/pipelane/pipelaner/internal/health"
	"github.com/pipelane/pipelaner/internal/metrics"
	"golang.org/x/sync/errgroup"
)

type Agent struct {
	pipelaner *Pipelaner

	hc      *health.Server
	metrics *metrics.Server
}

func NewAgent(file string) (*Agent, error) {
	ctx := context.Background()
	cfg, err := config.LoadFromPath(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	a := &Agent{}

	inits := []func(cfg *config.Pipelaner) error{
		a.initHealthCheck,
		a.initMetricsServer,
		a.initPipelaner,
	}

	for _, init := range inits {
		if err := init(cfg); err != nil {
			return nil, err
		}
	}

	return a, nil
}

func (a *Agent) initHealthCheck(cfg *config.Pipelaner) error {
	hcCfg := cfg.Settings.HealthCheck
	if hcCfg.Enable {
		hc, err := health.NewHealthCheck(hcCfg)
		if err != nil {
			return fmt.Errorf("init health check: %w", err)
		}
		a.hc = hc
	}
	return nil
}

func (a *Agent) initMetricsServer(cfg *config.Pipelaner) error {
	metricsCfg := cfg.Settings.Metrics
	if metricsCfg.Enable {
		m, err := metrics.NewMetricsServer(metricsCfg)
		if err != nil {
			return fmt.Errorf("init metrics server: %w", err)
		}
		a.metrics = m
	}
	return nil
}

func (a *Agent) initPipelaner(cfg *config.Pipelaner) error {
	pipelanerCfg := cfg.Pipelines
	logCfg := cfg.Settings.Logger
	// todo: use another solution for specific parameters
	p, err := NewPipelaner(
		pipelanerCfg,
		logCfg,
		cfg.Settings.Metrics.Enable,
		cfg.Settings.StartGCAfterMessageProcess,
	)
	if err != nil {
		return fmt.Errorf("init pipeliner: %w", err)
	}
	a.pipelaner = p
	return nil
}

func (a *Agent) Serve(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	if a.hc != nil {
		g.Go(func() error {
			return a.hc.Serve(ctx)
		})
	}
	if a.metrics != nil {
		g.Go(func() error {
			return a.metrics.Serve(ctx)
		})
	}
	g.Go(func() error {
		return a.pipelaner.Run(ctx)
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("run agent: %w", err)
	}
	return nil
}
