/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pipelane/pipelaner"
	"github.com/pipelane/pipelaner/example/pkl/gen/custom"
	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/pipeline/source"
)

// ============== Test generator ===============

type GenInt struct {
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

func (g *GenInt) Generate(_ context.Context, input chan<- any) {
	for {
		for i := 0; i < g.count; i++ {
			input <- i
		}
		time.Sleep(3 * time.Second)
	}
}

// ============= Test transform ===============

type TransMul struct {
	mul int
}

func (t *TransMul) Init(cfg transform.Transform) error {
	tCfg, ok := cfg.(custom.ExampleMul)
	if !ok {
		return errors.New("transform.Mul expects transform.TransMul")
	}
	t.mul = tCfg.GetMul()
	return nil
}

func (t *TransMul) Transform(val any) any {
	return t.mul * val.(int)
}

// ============= Test sink ==================

type Console struct {
}

func (c *Console) Init(cfg sink.Sink) error {
	return nil
}

func (c *Console) Sink(val any) {
	log.Println(val)
}

func init() {
	source.RegisterInput("example-generator", &GenInt{})
	source.RegisterTransform("example-mul", &TransMul{})
	source.RegisterSink("example-console", &Console{})
}

func main() {
	
	agent, err := pipelaner.NewAgent(
		"example/pkl/config.pkl",
	)
	if err != nil {
		panic(err)
	}
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	go func() {
		<-ctx.Done()
		os.Exit(100)
	}()
	defer stop()
	if err := agent.Serve(ctx); err != nil {
		panic(err)
	}
}
