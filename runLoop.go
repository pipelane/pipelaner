package pipelane

import (
	"context"
	"sort"
	"sync"
)

type MethodMap func(ctx context.Context, val any) any
type MethodSink func(ctx context.Context, val any)
type MethodGenerator func(ctx context.Context) any

type runLoop struct {
	*sync.RWMutex
	ctx          context.Context
	bufferSize   int64
	threadsCount *int64
	input        chan any
	outputs      []chan any
	transform    MethodMap
	sink         MethodSink
	receive      MethodGenerator
}

func (s *runLoop) SetMap(transform MethodMap) {
	s.transform = transform
}

func (s *runLoop) SetSink(sink MethodSink) {
	s.sink = sink
}

func (s *runLoop) SetGenerator(g MethodGenerator) {
	s.receive = g
}

func newRunLoop(
	ctx context.Context,
	bufferSize int64,
	threadsCount *int64,
) *runLoop {
	s := &runLoop{
		RWMutex:      &sync.RWMutex{},
		ctx:          ctx,
		bufferSize:   bufferSize,
		threadsCount: threadsCount,
		input:        make(chan any, bufferSize),
	}
	return s
}

func (s *runLoop) setInputChannel(ch chan any) {
	if s.input != nil {
		close(s.input)
	}
	s.input = ch
}

func (s *runLoop) Receive() {
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				break
			default:
				val := s.receive(s.ctx)
				s.input <- val
			}
		}
	}()
}

func (s *runLoop) run() {
	var sema chan struct{}
	if s.threadsCount != nil {
		sema = make(chan struct{}, *s.threadsCount)
	}
	semaphoreLockLock := func() {
		if sema != nil {
			sema <- struct{}{}
		}
	}
	semaphoreUnlock := func() {
		if sema != nil {
			<-sema
		}
	}
	go func() {
		for {
			select {
			case msg, ok := <-s.input:
				if !ok {
					return
				}
				s.rebalanced()
				semaphoreLockLock()
				go func(m any) {
					if s.transform != nil {
						m = s.transform(s.ctx, m)
					}
					s.RWMutex.RLock()
					for _, c := range s.outputs {
						c <- m
					}
					s.RWMutex.RUnlock()
					if s.sink != nil {
						s.sink(s.ctx, m)
					}
					semaphoreUnlock()
				}(msg)
			case <-s.ctx.Done():
				s.RWMutex.RLock()
				for _, c := range s.outputs {
					close(c)
				}
				s.RWMutex.RUnlock()
				return
			}
		}
	}()
}

func (s *runLoop) rebalanced() {
	s.RWMutex.Lock()
	sort.SliceIsSorted(s.outputs, func(i, j int) bool {
		diff1 := cap(s.outputs[i]) - len(s.outputs[i])
		diff2 := cap(s.outputs[j]) - len(s.outputs[j])
		return diff1 > diff2
	})
	s.RWMutex.Unlock()
}

func (s *runLoop) createOutput(bufferSize int64) chan any {
	ch := make(chan any, bufferSize)
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	s.outputs = append(s.outputs, ch)
	return ch
}
