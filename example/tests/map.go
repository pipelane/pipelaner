/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package tests

import (
	"context"
	"crypto/rand"
	"math/big"

	"pipelaner"
)

type MapGenerator struct {
}

func (i *MapGenerator) Init(cfg *pipelaner.BaseLaneConfig) error {
	return nil
}

func (i *MapGenerator) Generate(ctx context.Context, input chan<- any) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			input <- map[string]any{
				"user": map[string]any{
					"id": func() int64 {
						n, err := rand.Int(rand.Reader, big.NewInt(10))
						if err != nil {
							panic(err)
						}
						return n.Int64()
					}(),
				},
			}
		}
	}
}
