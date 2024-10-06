/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Agent struct {
	tree   *TreeLanes
	ctx    context.Context
	hc     *HealthCheck
	cancel context.CancelFunc
	cfg    *Config
}

func NewAgent(
	file string,
) (*Agent, error) {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	cfg, err := NewConfigFromFile(file)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	hc, err := NewHealthCheck(cfg.healthCheckConfig)
	if err != nil {
		return nil, fmt.Errorf("init healthcheck: %w", err)
	}
	return &Agent{
		tree:   nil,
		ctx:    ctx,
		hc:     hc,
		cancel: stop,
		cfg:    cfg,
	}, err
}

func (a *Agent) Serve() {
	if a.hc != nil {
		a.hc.Serve()
	}
	go func() {
		err := StartMetricsServer(a.cfg.metricsConfig)
		if err != nil {
			panic(err)
		}
	}()
	t, err := NewTreeFromConfig(a.ctx, a.cfg)
	if err != nil {
		panic(fmt.Errorf("init tree from config: %w", err))
	}
	a.tree = t
	<-a.ctx.Done()
}

func (a *Agent) Stop() {
	a.cancel()
}
