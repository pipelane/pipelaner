package pipelane

import (
	"context"
	"sort"
	"sync"
)

type transform func(ctx context.Context, val any) any
type sink func(ctx context.Context, val any)
type gen func(ctx context.Context) any

type subscriber struct {
	sync.RWMutex
	ctx        context.Context
	bufferSize int64
	input      chan any
	outputs    []chan any
	transform  transform
	sink       sink
	generate   gen
}

func (s *subscriber) Transform(transform transform) {
	s.transform = transform
}

func (s *subscriber) Sink(sink sink) {
	s.sink = sink
}

func (s *subscriber) Gen(g gen) {
	s.generate = g
}

func newSubscriber(
	ctx context.Context,
	bufferSize int64,
) *subscriber {
	s := &subscriber{
		RWMutex:    sync.RWMutex{},
		ctx:        ctx,
		bufferSize: bufferSize,
		input:      make(chan any, bufferSize),
	}
	return s
}

func (s *subscriber) setInputChannel(ch chan any) {
	if s.input != nil {
		close(s.input)
	}
	s.input = ch
}

func (s *subscriber) Generate() {
	for {
		select {
		case <-s.ctx.Done():
			break
		default:
			val := s.generate(s.ctx)
			s.input <- val
		}
	}
}

func (s *subscriber) run() {
	for {
		select {
		case msg, ok := <-s.input:
			if !ok {
				return
			}
			s.RWMutex.Lock()
			sort.SliceIsSorted(s.outputs, func(i, j int) bool {
				return len(s.outputs[i]) < len(s.outputs[j])
			})
			s.RWMutex.Unlock()
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
}

func (s *subscriber) createOutput(bufferSize int64) chan any {
	ch := make(chan any, bufferSize)
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	s.outputs = append(s.outputs, ch)
	return ch
}
