/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package components

import (
	"context"

	"github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/rs/zerolog"
)

type Logging interface {
	SetLogger(logger zerolog.Logger)
	Log() *zerolog.Logger
}

type Input interface {
	Init(cfg input.Input) error
	Generate(ctx context.Context, input chan<- any)
}

type Transform interface {
	Init(cfg transform.Transform) error
	Transform(val any) any
}
type Sequencer interface {
	Init(cfg transform.Sequencer) error
}

type Sink interface {
	Init(cfg sink.Sink) error
	Sink(val any) error
}

type TypeChecker interface {
	IsValidType(val any) bool
}

type Logger struct {
	logger *zerolog.Logger
}

func (l *Logger) SetLogger(logger zerolog.Logger) {
	l.logger = &logger
}

func (l *Logger) Log() *zerolog.Logger {
	return l.logger
}
