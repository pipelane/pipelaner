package node

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"github.com/LastPossum/kamino"
	configtransform "github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/internal/components"
	"github.com/pipelane/pipelaner/internal/metrics"
	"github.com/pipelane/pipelaner/internal/pipeline/source"
	"github.com/pipelane/pipelaner/internal/utils"
	"github.com/rs/zerolog"
)

const (
	transformNodeType = "transform"
)

type transformNodeCfg struct {
	name          string
	inputs        []string
	threadsCount  int
	outBufferSize int

	*nodeCfg
}

type Transform struct {
	impl          components.Transform
	cfg           *transformNodeCfg
	inputChannels []chan any
	outChannels   []chan any

	logger *zerolog.Logger
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
	if err := transformImpl.Init(cfg); err != nil {
		return nil, fmt.Errorf("init transform implementation: %w", err)
	}

	log := logger.With().
		Str("source", cfg.GetSourceName()).
		Str("type", transformNodeType).
		Str("lane_name", cfg.GetName()).
		Logger()
	return &Transform{
		impl: transformImpl,
		cfg: &transformNodeCfg{
			name:          cfg.GetName(),
			threadsCount:  cfg.GetThreads(),
			outBufferSize: cfg.GetOutputBufferSize(),
			inputs:        cfg.GetInputs(),
			nodeCfg:       buildOptions(opts...),
		},
		logger: &log,
	}, nil
}

func (t *Transform) AddInputChannel(ch chan any) {
	t.inputChannels = append(t.inputChannels, ch)
}

func (t *Transform) AddOutputChannel(ch chan any) {
	t.outChannels = append(t.outChannels, ch)
}

func (t *Transform) GetInputs() []string {
	return t.cfg.inputs
}

func (t *Transform) GetName() string {
	return t.cfg.name
}

func (t *Transform) GetOutputBufferSize() int {
	return t.cfg.outBufferSize
}

// Run non-blocking call that start Transform node action in separated goroutine
func (t *Transform) Run() error {
	if len(t.inputChannels) == 0 {
		return fmt.Errorf("no input channels configured for '%s'", t.cfg.name)
	}
	if len(t.outChannels) == 0 {
		return fmt.Errorf("no output channels configured for '%s'", t.cfg.name)
	}

	sema := utils.NewSemaphore(t.cfg.threadsCount)
	inChannel := utils.MergeChannels(t.inputChannels)

	go func() {
		var wg sync.WaitGroup

		for msg := range inChannel {
			wg.Add(1)
			sema.Acquire()
			go func() {
				defer wg.Done()
				defer sema.Release()

				msg = t.impl.Transform(msg)

				for _, ch := range t.outChannels {
					msg, err := t.prepareMessage(msg)
					if err != nil {
						t.logger.Error().Err(err).Msg("prepare message to send")
						continue
					}

					if err := t.preSendMessageAction(len(ch), cap(ch)); err != nil {
						t.logger.Error().Err(err).Msg("pre-send message action")
					}
					// send message
					ch <- msg
				}
			}()
		}
		t.logger.Debug().Msg("input channels processed")
		// Wait until all Transform goroutines finished
		wg.Wait()
		for _, ch := range t.outChannels {
			close(ch)
		}
	}()
	return nil
}

func (t *Transform) prepareMessage(msg any) (any, error) {
	if msg == nil {
		return nil, fmt.Errorf("received nil message")
	}
	switch m := msg.(type) {
	case error:
		if t.cfg.enableMetrics {
			metrics.TotalTransformationError.WithLabelValues(transformNodeType, t.cfg.name).Inc()
		}
		return nil, fmt.Errorf("received error: %w", m)
	default:
		kind := reflect.TypeOf(msg).Kind()
		switch kind {
		case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Struct:
			msg, err := kamino.Clone(msg)
			if err != nil {
				return nil, err
			}
			return msg, nil
		case reflect.Chan:
			// todo: broadcast channel
			return nil, nil
		default:
			return msg, nil
		}
	}
}

func (t *Transform) preSendMessageAction(length, capacity int) error {
	if t.cfg.enableMetrics {
		metrics.TotalMessagesCount.WithLabelValues(transformNodeType, t.cfg.name)
		metrics.BufferLength.WithLabelValues(transformNodeType, t.cfg.name).Set(float64(length))
		metrics.BufferCapacity.WithLabelValues(transformNodeType, t.cfg.name).Set(float64(capacity))
	}
	if t.cfg.callGC {
		runtime.GC()
	}
	return nil
}