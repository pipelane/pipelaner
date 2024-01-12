/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package sink

import (
	"fmt"

	"github.com/rs/zerolog"

	"pipelaner"
)

type Console struct {
	logger zerolog.Logger
}

func NewConsole(logger zerolog.Logger) *Console {
	return &Console{logger: logger}
}
func (c *Console) Init(ctx *pipelaner.Context) error {
	return nil
}

func (c *Console) Sink(ctx *pipelaner.Context, val any) {
	c.logger.Info().Msg(fmt.Sprintf("%v", val))
}
