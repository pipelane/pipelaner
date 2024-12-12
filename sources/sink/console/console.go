/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package console

import (
	"fmt"

	config "github.com/pipelane/pipelaner/gen/settings/logger"
	"github.com/pipelane/pipelaner/gen/settings/logger/loglevel"
	"github.com/pipelane/pipelaner/gen/source/sink"
	logger "github.com/pipelane/pipelaner/internal/logger"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/rs/zerolog"
)

func init() {
	source.RegisterSink("console", &Console{})
}

type Console struct {
	l *zerolog.Logger
}

func (c *Console) Init(cfg sink.Sink) error {
	cCfg, ok := cfg.(sink.Console)
	if !ok {
		return fmt.Errorf("invalid console config type: %T", cfg)
	}
	lCfg := config.Config{
		LogLevel:      loglevel.Info,
		EnableConsole: true,
		LogFormat:     cCfg.GetLogFormat(),
	}
	l, err := logger.NewLoggerWithCfg(&lCfg)
	if err != nil {
		return err
	}
	logs := l.With().
		Str("source", cfg.GetSourceName()).
		Str("type", "sink").
		Str("lane_name", cfg.GetName()).
		Logger()
	c.l = &logs
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
		c.l.Info().Msg(fmt.Sprintf("%v", val))
	}
}
