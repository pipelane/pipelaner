/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"sync"
)

type MethodMap func(ctx context.Context, val any) any
type MethodSink func(ctx context.Context, val any)
type MethodGenerator func(ctx context.Context, input chan<- any)

type loopCfg struct {
	bufferSize   int64
	threadsCount *int64
}

type methods struct {
	transform MethodMap
	sink      MethodSink
	generator MethodGenerator
}

type runLoop struct {
	mx            sync.RWMutex
	cfg           *loopCfg
	inputs        []chan any
	overrideInput bool
	outputs       []chan any
	methods       methods
}

func (s *runLoop) SetMap(transform MethodMap) {
	s.methods = methods{
		transform: transform,
	}
}

func (s *runLoop) SetSink(sink MethodSink) {
	s.methods = methods{
		sink: sink,
	}
}

func (s *runLoop) SetGenerator(g MethodGenerator) {
	s.methods = methods{
		generator: g,
	}
}

func newRunLoop(
	bufferSize int64,
	threadsCount *int64,
) *runLoop {
	s := &runLoop{
		mx: sync.RWMutex{},
		cfg: &loopCfg{
			bufferSize:   bufferSize,
			threadsCount: threadsCount,
		},
		inputs: []chan any{make(chan any, bufferSize)},
	}
	return s
}

func (s *runLoop) Receive(ctx context.Context) {
	for i := range s.inputs {
		go s.methods.generator(ctx, s.inputs[i])
	}
}

func (s *runLoop) run(ctx context.Context) {
	var sema chan struct{}
	if s.cfg.threadsCount != nil {
		sema = make(chan struct{}, *s.cfg.threadsCount)
	}
	semaphoreLock := func() {
		if sema != nil {
			sema <- struct{}{}
		}
	}
	semaphoreUnlock := func() {
		if sema != nil {
			<-sema
		}
	}
	closeSema := func() {
		if sema != nil {
			close(sema)
		}
	}
	input := mergeInputs(ctx, s.inputs...)
	go func() {
		defer closeSema()
		defer s.stop()
		for {
			select {
			case msg, ok := <-input:
				if !ok {
					return
				}
				if msg == nil {
					continue
				}
				semaphoreLock()
				go func(m any) {
					defer semaphoreUnlock()
					if s.methods.transform != nil {
						m = s.methods.transform(ctx, m)
					}
					if _, isErr := m.(error); isErr {
						return
					}
					if m == nil {
						return
					}
					for _, c := range s.outputs {
						c <- m
					}
					if s.methods.sink != nil {
						s.methods.sink(ctx, m)
					}

				}(msg)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *runLoop) createOutput(bufferSize int64) chan any {
	ch := make(chan any, bufferSize)
	s.mx.Lock()
	defer s.mx.Unlock()
	s.outputs = append(s.outputs, ch)
	return ch
}

func (s *runLoop) setInputChannel(ch chan any) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if !s.overrideInput {
		for i := range s.inputs {
			close(s.inputs[i])
		}
		s.inputs = []chan any{}
	}
	s.overrideInput = true
	s.inputs = append(s.inputs, ch)
}

func (s *runLoop) stop() {
	s.mx.Lock()
	defer s.mx.Unlock()
	if !s.overrideInput {
		for i := range s.inputs {
			close(s.inputs[i])
		}
	}
	for i := range s.outputs {
		close(s.outputs[i])
	}
	s.outputs = nil
}
