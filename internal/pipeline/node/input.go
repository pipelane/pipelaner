package node

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"

	"github.com/LastPossum/kamino"
	configinput "github.com/pipelane/pipelaner/gen/source/input"
	"github.com/pipelane/pipelaner/internal/components"
	"github.com/pipelane/pipelaner/internal/metrics"
	"github.com/pipelane/pipelaner/internal/pipeline/source"
	"github.com/rs/zerolog"
)

const (
	inputNodeType = "input"
)

type inputNodeCfg struct {
	name          string
	outBufferSize int

	metricsEnabled bool
	startGC        bool
}

type Input struct {
	impl           components.Input
	outputChannels []chan any
	cfg            *inputNodeCfg

	logger zerolog.Logger
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
	if err := inputImpl.Init(cfg); err != nil {
		return nil, fmt.Errorf("init input implementation: %s: %w", cfg.GetName(), err)
	}

	l := logger.With().
		Str("source", cfg.GetSourceName()).
		Str("type", inputNodeType).
		Str("lane_name", cfg.GetName()).
		Logger()
	return &Input{
		impl: inputImpl,
		cfg: &inputNodeCfg{
			name:           cfg.GetName(),
			outBufferSize:  cfg.GetOutputBufferSize(),
			metricsEnabled: false,
		},
		logger: l,
	}, nil
}

func (i *Input) AddOutputChannel(ch chan any) {
	i.outputChannels = append(i.outputChannels, ch)
}

func (i *Input) GetName() string {
	return i.cfg.name
}

func (i *Input) GetOutputBufferSize() int {
	return i.cfg.outBufferSize
}

func (i *Input) Run(ctx context.Context) error {
	if len(i.outputChannels) == 0 {
		return errors.New("no output channels configured")
	}

	input := make(chan any, i.cfg.outBufferSize*len(i.outputChannels))
	go func() {
		// close output channels on exit
		defer func() {
			for _, channel := range i.outputChannels {
				close(channel)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-input:
				if !ok {
					i.logger.Debug().Msg("input channel closed")
					return
				}
				for _, channel := range i.outputChannels {
					msg, err := i.prepareMessage(msg)
					if err != nil {
						i.logger.Error().Err(err).Msg("prepare message to send")
						continue
					}

					if err := i.preSendMessageAction(len(channel), cap(channel)); err != nil {
						i.logger.Error().Err(err).Msg("pre-send message action")
						continue
					}

					channel <- msg
				}

			}
		}
	}()

	go func() {
		defer close(input)
		// start generate messages
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
		//todo: так в итоге и не понял как должны обрабатываться ошибки
		// так как в run_loop ошибки input'a (generator'a не обрабатываются и прокидываются в transform)
		// я решил просто логгировать ошибку и не прокидывать ее дальше
		// - Предложение -
		// 1) сделать флаг, которым регулируется проброс ошибок.
		//   Флаг может конфигурироваться, как на уровне компонента (ноды),
		//   так и на уровне всего пайплайна
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

func (i *Input) preSendMessageAction(length, capacity int) error {
	//todo: не особо мне понятен как таковой смысл BufferLength и BufferCapacity метрик
	// особенно неизменяемое значение cap
	// + вопрос как данная метрика будет жить с несколькими выходными каналами для одной ноды?
	if i.cfg.metricsEnabled {
		metrics.TotalMessagesCount.WithLabelValues(inputNodeType, i.cfg.name)
		metrics.BufferLength.WithLabelValues(inputNodeType, i.cfg.name).Set(float64(length))
		metrics.BufferCapacity.WithLabelValues(inputNodeType, i.cfg.name).Set(float64(capacity))
	}
	//todo: в run_loop есть функционал, который вызывает runtime.GC(), перед отправкой сообщения
	// оставил здесь тоже, но не знаю насколько это понадобиться
	if i.cfg.startGC {
		runtime.GC()
	}
	return nil
}
