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

type Sequencer struct {
	commonTransform
	impl components.Sequencer
}

func NewSequencer(
	cfg configtransform.Transform,
	logger *zerolog.Logger,
	opts ...Option,
) (*Sequencer, error) {
	if cfg.GetName() == "" {
		return nil, errors.New("sequencer name is required")
	}
	if cfg.GetSourceName() == "" {
		return nil, errors.New("sequencer source name is required")
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

	transformImpl, err := source.GetSequencer(cfg.GetSourceName())

	if err != nil {
		return nil, fmt.Errorf("get sequencer implementation: %w", err)
	}
	l := logger.With().
		Str("source", cfg.GetSourceName()).
		Str("type", sequencerNodeType).
		Str("lane_name", cfg.GetName()).
		Logger()
	if v, ok := transformImpl.(components.Logging); ok {
		v.SetLogger(l)
	}

	return &Sequencer{
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
func (s *Sequencer) Run() error {
	if len(s.inputChannels) == 0 {
		return fmt.Errorf("no input channels configured for '%s'", s.cfg.name)
	}
	if len(s.outChannels) == 0 {
		return fmt.Errorf("no output channels configured for '%s'", s.cfg.name)
	}

	sema := synchronization.NewSemaphore(s.cfg.threadsCount)
	inChannel := synchronization.MergeInputs(s.inputChannels...)

	go func() {
		s.logger.Debug().Msg("starting sequencing messages")
		var wg sync.WaitGroup
		for msg := range inChannel {
			wg.Add(1)
			sema.Acquire()
			go func() {
				defer wg.Done()
				defer sema.Release()
				s.processingMessageByType(msg)
			}()
		}
		wg.Wait()
		for _, ch := range s.outChannels {
			close(ch)
		}
		s.logger.Debug().Msg("stop sequencer")
	}()
	return nil
}

func (s *Sequencer) processingMessageByType(msg any) {
	switch v := msg.(type) {
	case []any:
		for _, mV := range v {
			s.processingMessage(mV)
		}
	case chan any:
		for mV := range v {
			s.processingMessage(mV)
		}
	case AtomicData:
		s.processingAtomicMessage(v)
	default:
		s.processingMessage(v)
	}
}

func (s *Sequencer) processingMessage(msg any) {
	if e, ok := msg.(error); ok {
		if s.cfg.enableMetrics {
			metrics.TotalTransformationError.WithLabelValues(sequencerNodeType, s.cfg.name).Inc()
		}
		s.logger.Debug().Err(e).Msg("received error")
		return
	}
	for _, ch := range s.outChannels {
		mes, err := s.prepareMessage(msg)
		if err != nil {
			s.logger.Debug().Err(err).Msg("skip nil message sequencer")
			continue
		}
		s.preSendMessageAction(len(ch), cap(ch))
		ch <- mes
		s.postSinkAction()
	}
}

func (s *Sequencer) processingAtomicMessage(atomic any) {
	if e, ok := atomic.(error); ok {
		if s.cfg.enableMetrics {
			metrics.TotalTransformationError.WithLabelValues(sequencerNodeType, s.cfg.name).Inc()
		}
		s.logger.Debug().Err(e).Msg("received error")
		return
	}
	val, ok := atomic.(AtomicData)
	if !ok {
		s.logger.Debug().Err(errors.New("message is not atomic")).Msg("received error")
		return
	}
	switch data := val.Data().(type) {
	case []any:
		for _, mV := range data {
			s.atomicProcessSequence(mV, val)
		}
	case chan any:
		for mV := range data {
			s.atomicProcessSequence(mV, val)
		}
	default:
		s.atomicProcessSequence(val, val)
	}
}

func (s *Sequencer) atomicProcessSequence(mV any, val AtomicData) {
	for _, ch := range s.outChannels {
		mes, err := s.prepareMessage(mV)
		if err != nil {
			s.logger.Debug().Err(err).Msg("skip nil message sequencer")
			continue
		}
		s.preSendMessageAction(len(ch), cap(ch))
		newA := AtomicMessage{
			id:        val.ID(),
			data:      mes,
			successCh: val.Success(),
			errorCh:   val.Error(),
		}
		ch <- newA
		s.postSinkAction()
	}
}

func (s *Sequencer) preSendMessageAction(length, capacity int) {
	if s.cfg.enableMetrics {
		metrics.TotalMessagesCount.WithLabelValues(sequencerNodeType, s.cfg.name).Inc()
		metrics.BufferLength.WithLabelValues(sequencerNodeType, s.cfg.name).Set(float64(length))
		metrics.BufferCapacity.WithLabelValues(sequencerNodeType, s.cfg.name).Set(float64(capacity))
	}
}

func (s *Sequencer) postSinkAction() {
	if s.cfg.callGC {
		runtime.GC()
	}
}
