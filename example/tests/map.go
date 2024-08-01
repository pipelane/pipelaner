/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package tests

import (
	"crypto/rand"
	"math/big"

	"github.com/pipelane/pipelaner"
)

type MapGenerator struct {
}

func (i *MapGenerator) Init(_ *pipelaner.Context) error {
	return nil
}

func (i *MapGenerator) Generate(ctx *pipelaner.Context, input chan<- any) {
	for {
		select {
		case <-ctx.Context().Done():
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
