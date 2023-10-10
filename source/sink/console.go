package sink

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type Console struct {
	logger zerolog.Logger
}

func NewConsole(logger zerolog.Logger) *Console {
	return &Console{logger: logger}
}

func (c *Console) Sink(ctx context.Context, val any) {
	c.logger.Info().Msg(fmt.Sprintf("%v", val))
}
