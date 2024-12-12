/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package main

import (
	"context"
	"errors"
	"time"

	"github.com/pipelane/pipelaner"
	"github.com/pipelane/pipelaner/example/pkl/gen/custom"
	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/pipeline/components"
	"github.com/pipelane/pipelaner/pipeline/source"
	_ "github.com/pipelane/pipelaner/sources"
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
			i += 1
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

func init() {
	source.RegisterInput("example-generator", &GenInt{})
	source.RegisterTransform("example-mul", &TransMul{})
}

func main() {
	ctx := context.Background()
	agent, err := pipelaner.NewAgent(
		"example/pkl/config.pkl",
	)
	if err != nil {
		panic(err)
	}
	lock := make(chan struct{})
	go func() {
		time.Sleep(time.Second * 15)
		err = agent.Shutdown(context.Background())
		if err != nil {
			panic(err)
		}
		lock <- struct{}{}
	}()
	go func() {
		if err = agent.Serve(ctx); err != nil {
			panic(err)
		}
	}()
	<-lock
}
