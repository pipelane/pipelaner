package node

import (
	"errors"
	"fmt"
	"runtime"

	configsink "github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/internal/components"
	"github.com/pipelane/pipelaner/internal/metrics"
	"github.com/pipelane/pipelaner/internal/pipeline/source"
	"github.com/pipelane/pipelaner/internal/utils"
	"github.com/rs/zerolog"
)

const (
	sinkNodeType = "sink"
)

type sinkNodeCfg struct {
	name         string
	inputs       []string
	threadsCount int

	*nodeCfg
}

type Sink struct {
	impl          components.Sink
	cfg           *sinkNodeCfg
	inputChannels []chan any

	logger *zerolog.Logger
}

func NewSink(
	cfg configsink.Sink,
	logger *zerolog.Logger,
	opts ...Option,
) (*Sink, error) {
	if cfg.GetName() == "" {
		return nil, errors.New("must specify sink name")
	}
	if len(cfg.GetInputs()) == 0 {
		return nil, fmt.Errorf("'%s' has no inputs", cfg.GetName())
	}
	if cfg.GetThreads() < 1 {
		return nil, fmt.Errorf("'%s' invalid number of threads %d", cfg.GetName(), cfg.GetThreads())
	}
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	sinkImpl, err := source.GetSink(cfg.GetSourceName())
	if err != nil {
		return nil, fmt.Errorf("get sink implementation: %w", err)
	}
	if err := sinkImpl.Init(cfg); err != nil {
		return nil, fmt.Errorf("init sink implementation: %w", err)
	}

	log := logger.With().
		Str("source", cfg.GetSourceName()).
		Str("type", sinkNodeType).
		Str("lane_name", cfg.GetName()).
		Logger()
	return &Sink{
		impl: sinkImpl,
		cfg: &sinkNodeCfg{
			name:         cfg.GetName(),
			inputs:       cfg.GetInputs(),
			threadsCount: cfg.GetThreads(),
			nodeCfg:      buildOptions(opts...),
		},
		logger: &log,
	}, nil
}

func (s *Sink) AddInputChannel(inputChannel chan any) {
	s.inputChannels = append(s.inputChannels, inputChannel)
}

func (s *Sink) GetInputs() []string {
	return s.cfg.inputs
}

func (s *Sink) Run() error {
	if len(s.inputChannels) == 0 {
		return errors.New("no input channels configured")
	}

	inChannel := utils.MergeChannels(s.inputChannels)
	sema := utils.NewSemaphore(s.cfg.threadsCount)

	go func() {
		for msg := range inChannel {
			// process message in separated goroutine
			sema.Acquire()
			go func() {
				defer sema.Release()

				if err := s.preSinkAction(len(inChannel), cap(inChannel)); err != nil {
					s.logger.Error().Err(err).Msg("pre-sink action")
					return
				}
				// может имеет смысл возвращать ошибку?
				// для того чтобы логировать ее на уровне node + добавлять значение для метрик
				s.impl.Sink(msg)

				if err := s.postSinkAction(); err != nil {
					s.logger.Error().Err(err).Msg("post-sink action")
					return
				}
			}()
		}
		s.logger.Debug().Msg("input channels processed")
	}()
	return nil
}

// preSinkAction по логике из run_loop pipelaner'a
func (s *Sink) preSinkAction(length, capacity int) error {
	if s.cfg.enableMetrics {
		metrics.BufferLength.WithLabelValues(sinkNodeType, s.cfg.name).Set(float64(length))
		metrics.BufferCapacity.WithLabelValues(sinkNodeType, s.cfg.name).Set(float64(capacity))
	}
	if s.cfg.callGC {
		runtime.GC()
	}
	return nil
}

func (s *Sink) postSinkAction() error {
	if s.cfg.callGC {
		metrics.TotalMessagesCount.WithLabelValues(sinkNodeType, s.cfg.name).Inc()
	}
	return nil
}
