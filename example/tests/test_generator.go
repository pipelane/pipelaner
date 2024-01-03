/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package tests

import (
	"context"
	"time"

	"pipelaner"
)

type IntGenerator struct {
	inc uint64
}

func (i *IntGenerator) Init(cfg *pipelaner.BaseLaneConfig) error {
	i.inc = 1
	return nil
}
func (i *IntGenerator) Generate(ctx context.Context, input chan<- any) {
	for {
		i.inc += 1
		select {
		case <-ctx.Done():
			break
		default:
			input <- i.inc
		}
	}
}

type IntTwoGenerator struct {
}

func (i *IntTwoGenerator) Init(cfg *pipelaner.BaseLaneConfig) error {
	return nil
}

func (i *IntTwoGenerator) Generate(ctx context.Context, input chan<- any) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			input <- 2
		}
	}
}

type IntTransform struct {
}

func (i *IntTransform) New() pipelaner.Map {
	return &IntTransform{}
}

func (i *IntTransform) Map(ctx context.Context, val any) any {
	time.Sleep(time.Second)
	return val.(uint64) + 2
}

func (i *IntTransform) Init(cfg *pipelaner.BaseLaneConfig) error {
	return nil
}

type IntTransformEmpty struct {
}

func (i *IntTransformEmpty) New() pipelaner.Map {
	return &IntTransformEmpty{}
}

func (i *IntTransformEmpty) Map(ctx context.Context, val any) any {
	time.Sleep(time.Second)
	return val
}

func (i *IntTransformEmpty) Init(cfg *pipelaner.BaseLaneConfig) error {
	return nil
}
