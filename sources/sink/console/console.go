/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package console

import (
	"fmt"

	config "github.com/pipelane/pipelaner/gen/settings/logger"
	"github.com/pipelane/pipelaner/gen/settings/logger/loglevel"
	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/internal/logger"
	"github.com/pipelane/pipelaner/pipeline/node"
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
	l, err := logger.NewLoggerWithCfg(lCfg)
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

func (c *Console) Sink(val any) error {
	switch v := val.(type) {
	case node.AtomicData:
		err := c.Sink(v.Data())
		if err != nil {
			v.Error() <- v
			return err
		}
		v.Success() <- v
		return nil
	case chan node.AtomicData:
		for vals := range v {
			err := c.Sink(vals.Data())
			if err != nil {
				vals.Error() <- vals
				continue
			}
			vals.Success() <- vals
		}
		return nil
	case chan any:
		for vals := range v {
			err := c.Sink(vals)
			if err != nil {
				return err
			}
		}
		return nil
	case chan []byte:
		for vals := range v {
			err := c.Sink(vals)
			if err != nil {
				return err
			}
		}
		return nil
	case chan []string:
		for vals := range v {
			err := c.Sink(vals)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		c.l.Info().Msg(fmt.Sprintf("%v", val))
	}
	return nil
}
