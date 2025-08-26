/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package node

import (
	"errors"
	"fmt"
	"runtime"
	"sync"

	configtransform "github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/internal/metrics"
	"github.com/pipelane/pipelaner/internal/synchronization"
	"github.com/pipelane/pipelaner/pipeline/components"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/rs/zerolog"
)

type transformNodeCfg struct {
	name          string
	inputs        []string
	threadsCount  uint
	outBufferSize uint

	*nodeCfg
}

type Transform struct {
	commonTransform
	impl components.Transform
}

func NewTransform(
	cfg configtransform.Transform,
	logger *zerolog.Logger,
	opts ...Option,
) (*Transform, error) {
	if cfg.GetName() == "" {
		return nil, errors.New("transform name is required")
	}
	if cfg.GetSourceName() == "" {
		return nil, errors.New("transform source name is required")
	}
	if cfg.GetThreads() < 1 {
		return nil, fmt.Errorf("invalid number of threads %d", cfg.GetThreads())
	}
	if cfg.GetOutputBufferSize() < 1 {
		return nil, fmt.Errorf("invalid output buffer size %d", cfg.GetOutputBufferSize())
	}
	if len(cfg.GetInputs()) == 0 {
		return nil, errors.New("no input provided")
	}
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	transformImpl, err := source.GetTransform(cfg.GetSourceName())

	if err != nil {
		return nil, fmt.Errorf("get transform implementation: %w", err)
	}
	l := logger.With().
		Str("source", cfg.GetSourceName()).
		Str("type", transformNodeType).
		Str("lane_name", cfg.GetName()).
		Logger()
	if v, ok := transformImpl.(components.Logging); ok {
		v.SetLogger(l)
	}
	if err = transformImpl.Init(cfg); err != nil {
		return nil, fmt.Errorf("init transform implementation: %w", err)
	}

	return &Transform{
		impl: transformImpl,
		commonTransform: commonTransform{
			cfg: &transformNodeCfg{
				name:          cfg.GetName(),
				threadsCount:  cfg.GetThreads(),
				outBufferSize: cfg.GetOutputBufferSize(),
				inputs:        cfg.GetInputs(),
				nodeCfg:       buildOptions(opts...),
			},
			logger: l.With().Logger(),
		},
	}, nil
}

// Run non-blocking call that start Transform node action in separated goroutine.
func (t *Transform) Run() error {
	if len(t.inputChannels) == 0 {
		return fmt.Errorf("no input channels configured for '%s'", t.cfg.name)
	}
	if len(t.outChannels) == 0 {
		return fmt.Errorf("no output channels configured for '%s'", t.cfg.name)
	}

	sema := synchronization.NewSemaphore(t.cfg.threadsCount)
	inChannel := synchronization.MergeInputs(t.inputChannels...)

	go func() {
		t.logger.Debug().Msg("starting transform")
		var wg sync.WaitGroup

		for msg := range inChannel {
			wg.Add(1)
			sema.Acquire()
			go func() {
				var tmpMsg AtomicMessage
				defer wg.Done()
				defer sema.Release()
				msg = t.impl.Transform(msg)
				if e, ok := msg.(error); ok {
					if t.cfg.enableMetrics {
						metrics.TotalTransformationError.WithLabelValues(transformNodeType, t.cfg.name).Inc()
					}
					if m, oks := msg.(AtomicData); oks {
						if m.Error() != nil {
							m.Error() <- tmpMsg
						}
					}
					t.logger.Debug().Err(e).Msg("received error")
					return
				}
				for _, ch := range t.outChannels {
					var mes any
					var err error
					switch ms := msg.(type) {
					case AtomicData:
						mes, err = t.prepareMessage(ms.Data())
						mes = ms.UpdateData(mes)
					default:
						mes, err = t.prepareMessage(ms)
					}
					if err != nil {
						t.logger.Debug().Err(err).Msg("skip nil message transform")
						continue
					}
					t.preSendMessageAction(len(ch), cap(ch))
					ch <- mes
					t.postSinkAction()
				}
			}()
		}
		wg.Wait()
		for _, ch := range t.outChannels {
			close(ch)
		}
		t.logger.Debug().Msg("stop transform")
	}()
	return nil
}

func (t *Transform) preSendMessageAction(length, capacity int) {
	if t.cfg.enableMetrics {
		metrics.TotalMessagesCount.WithLabelValues(transformNodeType, t.cfg.name).Inc()
		metrics.BufferLength.WithLabelValues(transformNodeType, t.cfg.name).Set(float64(length))
		metrics.BufferCapacity.WithLabelValues(transformNodeType, t.cfg.name).Set(float64(capacity))
	}
}

func (t *Transform) postSinkAction() {
	if t.cfg.callGC {
		runtime.GC()
	}
}
