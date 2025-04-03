/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package node

import (
	"fmt"
	"reflect"

	"github.com/LastPossum/kamino"
	"github.com/rs/zerolog"
)

type Type string

type nodeCfg struct {
	enableMetrics bool
	callGC        bool
}

const (
	transformNodeType = "transform"
	sequencerNodeType = "sequencer"
)

type Option func(*nodeCfg)

func WithMetrics() Option {
	return func(cfg *nodeCfg) {
		cfg.enableMetrics = true
	}
}

func WithCallGC() Option {
	return func(cfg *nodeCfg) {
		cfg.callGC = true
	}
}

func buildOptions(opts ...Option) *nodeCfg {
	cfg := &nodeCfg{}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

type commonTransform struct {
	inputChannels []chan any
	outChannels   []chan any
	cfg           *transformNodeCfg
	logger        zerolog.Logger
}

func (t *commonTransform) AddInputChannel(ch chan any) {
	t.inputChannels = append(t.inputChannels, ch)
}

func (t *commonTransform) AddOutputChannel(ch chan any) {
	t.outChannels = append(t.outChannels, ch)
}

func (t *commonTransform) GetInputs() []string {
	return t.cfg.inputs
}

func (t *commonTransform) GetName() string {
	return t.cfg.name
}

func (t *commonTransform) GetOutputBufferSize() uint {
	return t.cfg.outBufferSize
}

func (t *commonTransform) prepareMessage(msg any) (any, error) {
	if msg == nil {
		return nil, fmt.Errorf("received nil message")
	}
	kind := reflect.TypeOf(msg).Kind()
	switch kind {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Struct:
		copyMsg, err := kamino.Clone(msg)
		if err != nil {
			return nil, err
		}
		return copyMsg, nil
	case reflect.Chan:
		if len(t.outChannels) == 1 {
			return msg, nil
		}
		mCh, ok := msg.(chan any)
		if !ok {
			return nil, fmt.Errorf("received message of type '%T'", msg)
		}
		newChan := make(chan any, cap(mCh))
		go func() {
			close(newChan)
			for m := range mCh {
				copyMes, errs := t.prepareMessage(m)
				if errs != nil {
					return
				}
				newChan <- copyMes
			}
		}()
		return newChan, nil
	default:
		return msg, nil
	}
}
