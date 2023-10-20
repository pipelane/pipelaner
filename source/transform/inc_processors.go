package transform

import (
	"context"
	"time"

	pipelane "github.com/pipelane/pipelaner"
)

type IncProcessor struct {
}

func (i *IncProcessor) Init(cfg *pipelane.BaseLaneConfig) error {
	return nil
}

func (i *IncProcessor) Map(ctx context.Context, val any) any {
	v := val.(int)
	v++
	time.Sleep(time.Second * 5)
	return v
}
