/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package main

import (
	"context"
	"errors"
	"time"

	"github.com/pipelane/pipelaner/examples/custom/gen/custom"
	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/pipeline/components"
)

// ============== Test generator ===============

type GenInt struct {
	components.Logger
	count int
}

func (g *GenInt) Init(cfg input.Input) error {
	gCfg, ok := cfg.(custom.ExampleGenInt)
	if !ok {
		return errors.New("invalid config")
	}
	g.count = gCfg.GetCount()
	return nil
}

func (g *GenInt) Generate(ctx context.Context, input chan<- any) {
	i := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
			input <- i
			i++
			time.Sleep(time.Second * 1)
		}
	}
}

// ============= Test transform ===============

type TransMul struct {
	mul int
}

func (t *TransMul) Init(cfg transform.Transform) error {
	tCfg, ok := cfg.(custom.ExampleMul)
	if !ok {
		return errors.New("invalid config")
	}
	t.mul = tCfg.GetMul()
	return nil
}

func (t *TransMul) Transform(val any) any {
	v, ok := val.(int)
	if !ok {
		return errors.New("invalid value")
	}
	return t.mul * v
}
