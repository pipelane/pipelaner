/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type Agent struct {
	tree   *TreeLanes
	ctx    context.Context
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
		return nil, err
	}
	t, err := NewTreeFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Agent{
		tree:   t,
		ctx:    ctx,
		cancel: stop,
	}, err
}

func (a *Agent) Serve() {
	<-a.ctx.Done()
}

func (a *Agent) Stop() {
	a.cancel()
}
