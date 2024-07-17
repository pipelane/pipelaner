/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"reflect"
	"sync"
	"sync/atomic"
)

type MethodMap func(ctx *Context, val any) any
type MethodSink func(ctx *Context, val any)
type MethodGenerator func(ctx *Context, input chan<- any)

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
	stopped       atomic.Bool
	outputs       []chan any
	methods       methods
	context       *Context
}

func (s *runLoop) setContext(context *Context) {
	s.context = context
}

func (s *runLoop) setMap(transform MethodMap) {
	s.methods = methods{
		transform: transform,
	}
}

func (s *runLoop) setSink(sink MethodSink) {
	s.methods = methods{
		sink: sink,
	}
}

func (s *runLoop) setGenerator(g MethodGenerator) {
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

func (s *runLoop) receive() {
	for i := range s.inputs {
		c := s.inputs[i]
		go func(ch chan any) {
			s.methods.generator(s.context, ch)
			close(ch)
		}(c)
	}
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
	input := mergeInputs(s.context.Context(), s.inputs...)
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
				valMsg := reflect.ValueOf(msg)
				if valMsg.Kind() == reflect.Pointer {
					msg = valMsg.Elem().Interface()
				}
				semaphoreLock()
				go func(m any) {
					defer semaphoreUnlock()
					if s.methods.transform != nil {
						m = s.methods.transform(s.context, m)
					}
					if _, isErr := m.(error); isErr {
						return
					}
					if m == nil {
						return
					}
					s.mx.RLock()
					for _, c := range s.outputs {
						c <- m
					}
					s.mx.RUnlock()
					if s.methods.sink != nil {
						s.methods.sink(s.context, m)
					}
				}(msg)
			case <-s.context.Context().Done():
				s.stopped.Store(true)
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
	for i := range s.outputs {
		close(s.outputs[i])
	}
	s.outputs = nil
}
