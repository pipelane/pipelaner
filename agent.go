/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"fmt"
	"time"

	config "github.com/pipelane/pipelaner/gen/pipelaner"
	"github.com/pipelane/pipelaner/internal/health"
	"github.com/pipelane/pipelaner/internal/metrics"
	"github.com/pipelane/pipelaner/internal/migrator"
	"github.com/pipelane/pipelaner/internal/pprof"
	"golang.org/x/sync/errgroup"
)

type Agent struct {
	pipelaner *Pipelaner

	hc      *health.Server
	metrics *metrics.Server
	pprof   *pprof.Server
	cancel  context.CancelFunc
	cfg     *config.Pipelaner
}

func NewAgent(file string) (*Agent, error) {
	ctx := context.Background()
	cfg, err := config.LoadFromPath(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	a := &Agent{
		cfg: &cfg,
	}

	inits := []func(cfg *config.Pipelaner) error{
		a.initAndMigrate,
		a.initHealthCheck,
		a.initMetricsServer,
		a.initPprofServer,
		a.initPipelaner,
	}

	for _, init := range inits {
		if err = init(&cfg); err != nil {
			return nil, err
		}
	}

	return a, nil
}

func (a *Agent) initAndMigrate(cfg *config.Pipelaner) error {
	mCfg := cfg.Settings.Migrations
	if mCfg != nil {
		m, err := migrator.NewMigrator(cfg)
		if err != nil {
			return fmt.Errorf("migration: %w", err)
		}
		err = m.Migrate()
		if err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	return nil
}

func (a *Agent) initHealthCheck(cfg *config.Pipelaner) error {
	hcCfg := cfg.Settings.HealthCheck
	if hcCfg != nil {
		hc, err := health.NewHealthCheck(cfg)
		if err != nil {
			return fmt.Errorf("init health check: %w", err)
		}
		a.hc = hc
	}
	return nil
}

func (a *Agent) initMetricsServer(cfg *config.Pipelaner) error {
	metricsCfg := cfg.Settings.Metrics
	if metricsCfg != nil {
		m, err := metrics.NewMetricsServer(metricsCfg)
		if err != nil {
			return fmt.Errorf("init metrics server: %w", err)
		}
		a.metrics = m
	}
	return nil
}

func (a *Agent) initPprofServer(cfg *config.Pipelaner) error {
	pprofCfg := cfg.Settings.Pprof
	if pprofCfg != nil {
		p := pprof.NewServer(pprofCfg)
		a.pprof = p
	}
	return nil
}

func (a *Agent) initPipelaner(cfg *config.Pipelaner) error {
	pipelanerCfg := cfg.Pipelines
	logCfg := cfg.Settings.Logger
	metricsCfg := cfg.Settings.Metrics
	mEnable := metricsCfg != nil
	p, err := NewPipelaner(
		pipelanerCfg,
		logCfg,
		mEnable,
		cfg.Settings.StartGCAfterMessageProcess,
	)
	if err != nil {
		return fmt.Errorf("init pipeliner: %w", err)
	}
	a.pipelaner = p
	return nil
}

func (a *Agent) Serve(ctx context.Context) error {
	g := errgroup.Group{}
	inputsCtx, cancel := context.WithCancel(ctx)
	a.cancel = cancel
	if a.hc != nil {
		g.Go(func() error {
			return a.hc.Serve(ctx)
		})
	}
	if a.pprof != nil {
		g.Go(func() error {
			return a.pprof.Serve(ctx)
		})
	}
	if a.metrics != nil {
		g.Go(func() error {
			return a.metrics.Serve(ctx)
		})
	}
	g.Go(func() error {
		return a.pipelaner.Run(inputsCtx)
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("run agent: %w", err)
	}
	return nil
}

func (a *Agent) Shutdown(ctx context.Context) error {
	a.cancel()
	time.Sleep(a.cfg.Settings.GracefulShutdownDelay.GoDuration())
	if a.pprof != nil {
		err := a.pprof.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("shutdown pprof: %w", err)
		}
	}
	if a.metrics != nil {
		err := a.metrics.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("shutdown metrics: %w", err)
		}
	}
	if a.hc != nil {
		err := a.hc.Shutdown()
		if err != nil {
			return fmt.Errorf("shutdown healthcheck: %w", err)
		}
	}
	return nil
}
