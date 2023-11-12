/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package generator

import (
	"context"

	pipelane "github.com/pipelane/pipelaner"
)

type IntGenerator struct {
}

func (i *IntGenerator) Init(cfg *pipelane.BaseLaneConfig) error {
	return nil
}

func (i *IntGenerator) Generate(ctx context.Context, input chan<- any) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			input <- 1
		}
	}
}
