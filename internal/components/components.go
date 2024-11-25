package components

import (
	"context"

	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/gen/source/transform"
)

type Input interface {
	Init(cfg input.Input) error
	Generate(ctx context.Context, input chan<- any)
}

type Transform interface {
	Init(cfg transform.Transform) error
	Transform(val any) any
}

type Sink interface {
	Init(cfg sink.Sink) error
	Sink(val any)
}
