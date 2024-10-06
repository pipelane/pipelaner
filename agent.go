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
	cfg     *Config
	tree    *TreeLanes
	hc      *HealthCheck
	metrics *MetricsServer
	ctx     context.Context
	cancel  context.CancelFunc
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

	hc, err := NewHealthCheck(*cfg)
	if err != nil {
		return nil, fmt.Errorf("init healthcheck server: %w", err)
	}
	m, err := NewMetricsServer(*cfg)
	if err != nil {
		return nil, fmt.Errorf("init metrics server: %w", err)
	}
	return &Agent{
		tree:    nil,
		ctx:     ctx,
		hc:      hc,
		metrics: m,
		cancel:  stop,
		cfg:     cfg,
	}, err
}

func (a *Agent) Serve() {
	if a.hc != nil {
		a.hc.Serve()
	}
	go func() {
		if a.metrics != nil {
			err := a.metrics.Serve()
			if err != nil {
				panic(err)
			}
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
