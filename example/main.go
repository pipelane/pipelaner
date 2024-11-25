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

	"github.com/pipelane/pipelaner"
	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/internal/pipeline/source"
	_ "github.com/pipelane/pipelaner/source"
)

// ============== Test generator ===============

type GenInt struct {
	a string
}

func (g *GenInt) Init(cfg input.Input) error {
	return nil
}

func (g *GenInt) Generate(_ context.Context, input chan<- any) {
	for i := 0; i < 10; i++ {
		input <- i
	}
}

// ============= Test transform ===============

type TransMul struct {
	mul int
}

func (t *TransMul) Init(cfg transform.Transform) error {
	tCfg, ok := cfg.(transform.Mul)
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

func main() {
	log.Println("Start")
	source.RegisterInput("test_gen", &GenInt{})
	source.RegisterTransform("mul", &TransMul{})
	source.RegisterSink("console", &Console{})
	agent, err := pipelaner.NewAgent(
		"/Users/n.frolov/GolandProjects/pipelaner_old/pipelaner_pkl/pipelaner/pkl/dev/config.pkl",
	)
	if err != nil {
		panic(err)
	}
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	agent.Serve(ctx)
}
