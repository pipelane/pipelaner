package transform

import (
	"context"
	"time"
)

type IncProcessor struct {
}

func (i IncProcessor) Transform(ctx context.Context, val any) any {
	v := val.(int)
	v++
	time.Sleep(time.Second * 5)
	return v
}
