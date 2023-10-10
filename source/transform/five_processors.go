package transform

import (
	"context"
	"time"
)

type FiveProcessor struct {
}

func (i FiveProcessor) Transform(ctx context.Context, val any) any {
	v := val.(int)
	v += 5
	time.Sleep(time.Second * 5)
	return v
}
