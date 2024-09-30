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
	t, err := NewTreeFromConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("init tree from config: %w", err)
	}

	hc, err := NewHealthCheck(cfg.healthCheckConfig)
	if err != nil {
		return nil, fmt.Errorf("init healthcheck: %w", err)
	}
	return &Agent{
		tree:   t,
		ctx:    ctx,
		hc:     hc,
		cancel: stop,
	}, err
}

func (a *Agent) Serve() {
	if a.hc != nil {
		a.hc.Serve()
	}

	<-a.ctx.Done()
}

func (a *Agent) Stop() {
	a.cancel()
}
