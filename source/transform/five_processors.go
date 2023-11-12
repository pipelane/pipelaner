/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package transform

import (
	"context"
	"time"

	pipelane "github.com/pipelane/pipelaner"
)

type FiveProcessor struct {
}

func (i *FiveProcessor) Init(cfg *pipelane.BaseLaneConfig) error {
	return nil
}

func (i *FiveProcessor) Map(ctx context.Context, val any) any {
	v := val.(int)
	v += 5
	time.Sleep(time.Second * 5)
	return v
}
