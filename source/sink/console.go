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

func init() {
	pipelaner.RegisterSink("console", NewConsole())
}

func NewConsole() *Console {
	return &Console{}
}
func (c *Console) Init(ctx *pipelaner.Context) error {
	return nil
}

func (c *Console) Sink(ctx *pipelaner.Context, val any) {
	c.logger.Info().Msg(fmt.Sprintf("%v", val))
}
