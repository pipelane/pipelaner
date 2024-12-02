/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package console

import (
	"fmt"

	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/internal/logger"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/rs/zerolog"
)

func init() {
	source.RegisterSink("console", &Console{})
}

type Console struct {
	logger *zerolog.Logger
}

func (c *Console) Init(cfg sink.Sink) error {
	_, ok := cfg.(sink.Console)
	if !ok {
		return fmt.Errorf("invalid console config type: %T", cfg)
	}

	l := logger.NewLogger()
	c.logger = &l
	return nil
}

func (c *Console) Sink(val any) {
	switch v := val.(type) {
	case chan any:
		for vals := range v {
			c.Sink(vals)
		}
		return
	case chan []byte:
		for vals := range v {
			c.Sink(vals)
		}
		return
	case chan []string:
		for vals := range v {
			c.Sink(vals)
		}
		return
	default:
		c.logger.Info().Msg(fmt.Sprintf("%v", val))
	}
}
