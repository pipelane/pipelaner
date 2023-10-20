package sink

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	pipelane "github.com/pipelane/pipelaner"
)

type Console struct {
	logger zerolog.Logger
}

func NewConsole(logger zerolog.Logger) *Console {
	return &Console{logger: logger}
}
func (c *Console) Init(cfg *pipelane.BaseLaneConfig) error {
	return nil
}

func (c *Console) Sink(ctx context.Context, val any) {
	c.logger.Info().Msg(fmt.Sprintf("%v", val))
}
