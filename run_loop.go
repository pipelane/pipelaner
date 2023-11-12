/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"sort"
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
	ctx           context.Context
	cfg           *loopCfg
	input         chan any
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
	ctx context.Context,
	bufferSize int64,
	threadsCount *int64,
) *runLoop {
	s := &runLoop{
		mx:  sync.RWMutex{},
		ctx: ctx,
		cfg: &loopCfg{
			bufferSize:   bufferSize,
			threadsCount: threadsCount,
		},
		input: make(chan any, bufferSize),
	}
	return s
}

func (s *runLoop) setInputChannel(ch chan any) {
	if s.input != nil {
		close(s.input)
	}
	s.overrideInput = true
	s.input = ch
}

func (s *runLoop) Receive() {
	go s.methods.generator(s.ctx, s.input)
}

func (s *runLoop) run() {
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
	go func() {
		defer closeSema()
		defer s.stop()
		for {
			select {
			case msg, ok := <-s.input:
				if !ok {
					return
				}
				s.rebalanced()
				semaphoreLock()
				go func(m any) {
					defer semaphoreUnlock()
					if s.methods.transform != nil {
						m = s.methods.transform(s.ctx, m)
					}
					if _, isErr := m.(error); isErr {
						return
					}
					if m != nil {
						s.mx.RLock()
						for _, c := range s.outputs {
							c <- m
						}
						s.mx.RUnlock()
					}
					if s.methods.sink != nil {
						s.methods.sink(s.ctx, m)
					}

				}(msg)
			case <-s.ctx.Done():
				return
			}
		}
	}()
}

func (s *runLoop) rebalanced() {
	s.mx.Lock()
	sort.SliceIsSorted(s.outputs, func(i, j int) bool {
		diff1 := cap(s.outputs[i]) - len(s.outputs[i])
		diff2 := cap(s.outputs[j]) - len(s.outputs[j])
		return diff1 > diff2
	})
	s.mx.Unlock()
}

func (s *runLoop) createOutput(bufferSize int64) chan any {
	ch := make(chan any, bufferSize)
	s.mx.Lock()
	defer s.mx.Unlock()
	s.outputs = append(s.outputs, ch)
	return ch
}

func (s *runLoop) stop() {
	s.mx.Lock()
	defer s.mx.Unlock()
	if !s.overrideInput {
		close(s.input)
	}
	for i := range s.outputs {
		close(s.outputs[i])
	}
	s.outputs = nil
}
