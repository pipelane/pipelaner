/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/LastPossum/kamino"
	"github.com/prometheus/client_golang/prometheus"
)

var totalMessagesCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "message_total",
		Help: "Total number of messages.",
	},
	[]string{"type", "name"},
)

var totalTransformationError = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "transformations_error_total",
		Help: "Total number of errors.",
	},
	[]string{"type", "name"},
)

var bufferCapacity = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "buffer_capacity",
		Help: "Buffer capacity.",
	},
	[]string{"type", "name"},
)

var bufferLength = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "buffer_length",
		Help: "Buffer length.",
	},
	[]string{"type", "name"},
)

func init() {
	prometheus.MustRegister(totalMessagesCount)
	prometheus.MustRegister(totalTransformationError)
	prometheus.MustRegister(bufferCapacity)
	prometheus.MustRegister(bufferLength)
}

type MethodMap func(ctx *Context, val any) any
type MethodSink func(ctx *Context, val any)
type MethodGenerator func(ctx *Context, input chan<- any)

type loopCfg struct {
	bufferSize   int64
	threadsCount int64
	startGC      bool
}

type methods struct {
	transform MethodMap
	sink      MethodSink
	generator MethodGenerator
}

type runLoop struct {
	cfg           *loopCfg
	stopped       atomic.Bool
	methods       methods
	context       *Context
	metricsEnable bool

	mx            sync.RWMutex
	overrideInput bool
	inputs        []chan any
	outputs       []chan any
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
	threadsCount int64,
	startGC bool,
	isMetricsEnabled bool,
) *runLoop {
	s := &runLoop{
		mx: sync.RWMutex{},
		cfg: &loopCfg{
			bufferSize:   bufferSize,
			threadsCount: threadsCount,
			startGC:      startGC,
		},
		inputs:        []chan any{make(chan any, bufferSize)},
		metricsEnable: isMetricsEnabled,
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

func (s *runLoop) start() {
	var sema chan struct{}
	if s.cfg.threadsCount != 0 {
		sema = make(chan struct{}, s.cfg.threadsCount)
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
				if s.metricsEnable && s.context.LaneType() == InputType {
					totalMessagesCount.
						WithLabelValues(
							string(s.context.LaneType()),
							s.context.LaneName(),
						).Inc()
				}
				if s.metricsEnable {
					bufferLength.WithLabelValues(
						string(s.context.LaneType()),
						s.context.LaneName()).Set(float64(len(input)))
					bufferCapacity.WithLabelValues(
						string(s.context.LaneType()),
						s.context.LaneName()).Set(float64(cap(input)))
				}
				if !ok {
					return
				}
				if msg == nil {
					continue
				}
				kind := reflect.TypeOf(msg).Kind()
				switch kind {
				case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Struct:
					m, err := kamino.Clone(msg)
					if err != nil {
						continue
					}
					msg = m
				default:
				}
				semaphoreLock()
				go s.produceMessages(semaphoreUnlock, msg)
			case <-s.context.Context().Done():
				s.stopped.Store(true)
				return
			}
			if s.cfg.startGC {
				runtime.GC()
			}
		}
	}()
}

func (s *runLoop) produceMessages(unlock func(), m any) {
	defer unlock()
	if s.methods.transform != nil {
		m = s.methods.transform(s.context, m)
	}
	_, isErr := m.(error)
	if isErr && s.metricsEnable {
		totalTransformationError.
			WithLabelValues(
				string(s.context.LaneType()),
				s.context.LaneName(),
			).
			Inc()
		logger := s.context.Logger()
		logger.Error().
			Err(m.(error)).Msg("run loop error")
		return
	}
	if m == nil {
		return
	}
	s.mx.RLock()
	// check message type before send to output
	switch mVal := m.(type) {
	case chan any: // temp contract solution
		broadcastChannels(s.outputs, mVal)
	default:
		for _, c := range s.outputs {
			c <- m
		}
	}
	s.mx.RUnlock()
	if s.methods.sink != nil {
		s.methods.sink(s.context, m)
	}
	if s.metricsEnable {
		totalMessagesCount.
			WithLabelValues(
				string(s.context.LaneType()),
				s.context.LaneName(),
			).Inc()
	}
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
