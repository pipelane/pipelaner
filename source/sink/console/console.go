/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package console

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/pipelane/pipelaner"
)

type Console struct {
	logger *zerolog.Logger
}

func init() {
	pipelaner.RegisterSink("console", &Console{})
}

func (c *Console) Init(ctx *pipelaner.Context) error {
	c.logger = ctx.Logger()
	return nil
}

func (c *Console) Sink(_ *pipelaner.Context, val any) {
	switch v := val.(type) {
	case chan any:
		for vals := range v {
			c.logger.Info().Msg(fmt.Sprintf("%v", vals))
		}
		return
	default:
		c.logger.Info().Msg(fmt.Sprintf("%v", val))
	}

}
