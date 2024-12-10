package node

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"

	"github.com/LastPossum/kamino"
	configinput "github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/internal/metrics"
	"github.com/pipelane/pipelaner/pipeline/components"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/rs/zerolog"
)

const (
	inputNodeType = "input"
)

type inputNodeCfg struct {
	name          string
	outBufferSize int

	*nodeCfg
}

type Input struct {
	impl        components.Input
	outChannels []chan any
	cfg         *inputNodeCfg
	logger      zerolog.Logger
}

func NewInput(cfg configinput.Input, logger *zerolog.Logger, opts ...Option) (*Input, error) {
	if cfg.GetName() == "" {
		return nil, errors.New("input name is required")
	}
	if cfg.GetSourceName() == "" {
		return nil, errors.New("input source name is required")
	}
	if cfg.GetThreads() < 1 {
		return nil, fmt.Errorf("invalid number of threads %d", cfg.GetThreads())
	}
	if cfg.GetOutputBufferSize() < 1 {
		return nil, fmt.Errorf("invalid output buffer size %d", cfg.GetOutputBufferSize())
	}
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	inputImpl, err := source.GetInput(cfg.GetSourceName())
	if err != nil {
		return nil, fmt.Errorf("get input implementation: %w", err)
	}
	l := logger.With().
		Str("source", cfg.GetSourceName()).
		Str("type", inputNodeType).
		Str("lane_name", cfg.GetName()).
		Logger()
	if v, ok := inputImpl.(components.Logging); ok {
		v.SetLogger(l)
	}
	if err = inputImpl.Init(cfg); err != nil {
		return nil, fmt.Errorf("init input implementation: %s: %w", cfg.GetName(), err)
	}
	return &Input{
		impl: inputImpl,
		cfg: &inputNodeCfg{
			name:          cfg.GetName(),
			outBufferSize: cfg.GetOutputBufferSize(),
			nodeCfg:       buildOptions(opts...),
		},
		logger: l.With().Logger(),
	}, nil
}

func (i *Input) AddOutputChannel(ch chan any) {
	i.outChannels = append(i.outChannels, ch)
}

func (i *Input) GetName() string {
	return i.cfg.name
}

func (i *Input) GetOutputBufferSize() int {
	return i.cfg.outBufferSize
}

func (i *Input) Run(ctx context.Context) error {
	if len(i.outChannels) == 0 {
		return errors.New("no output channels configured")
	}

	input := make(chan any, i.cfg.outBufferSize*len(i.outChannels))
	go func() {
		defer func() {
			for _, channel := range i.outChannels {
				close(channel)
			}
		}()

		for msg := range input {
			for _, ch := range i.outChannels {
				m, err := i.prepareMessage(msg)
				if err != nil {
					i.logger.Error().Err(err).Msg("prepare message failed")
					continue
				}
				i.preSendMessageAction(len(input), len(input))
				ch <- m
			}
		}
		i.logger.Debug().Msg("input channel closed")
	}()

	go func() {
		defer close(input)
		i.impl.Generate(ctx, input)
	}()
	return nil
}

func (i *Input) prepareMessage(msg any) (any, error) {
	if msg == nil {
		return nil, errors.New("received nil message")
	}
	switch m := msg.(type) {
	case error:
		if i.cfg.enableMetrics {
			metrics.TotalTransformationError.WithLabelValues(inputNodeType, i.cfg.name).Inc()
		}
		return nil, fmt.Errorf("received error: %w", m)
	default:
		kind := reflect.TypeOf(msg).Kind()
		switch kind {
		case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Struct:
			mes, err := kamino.Clone(msg)
			if err != nil {
				return nil, err
			}
			return mes, nil
		default:
			return msg, nil
		}
	}
}

func (i *Input) preSendMessageAction(length, capacity int) {
	if i.cfg.enableMetrics {
		metrics.TotalMessagesCount.WithLabelValues(inputNodeType, i.cfg.name).Inc()
		metrics.BufferLength.WithLabelValues(inputNodeType, i.cfg.name).Set(float64(length))
		metrics.BufferCapacity.WithLabelValues(inputNodeType, i.cfg.name).Set(float64(capacity))
	}
	if i.cfg.callGC {
		runtime.GC()
	}
}
