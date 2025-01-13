/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package batch

import (
	"fmt"
	"sync"

	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/pipeline/source"
)

func init() {
	source.RegisterTransform("batch", &Batch{})
}

type Config struct {
	Size int64 `pipelane:"size"`
}

type Batch struct {
	mx   sync.Mutex
	ch   chan any
	size uint
}

func (b *Batch) Init(cfg transform.Transform) error {
	bCfg, ok := cfg.(transform.Batch)
	if !ok {
		return fmt.Errorf("invalid batch config type: %T", cfg)
	}
	b.size = bCfg.GetSize()
	b.ch = make(chan any, bCfg.GetSize())
	return nil
}

func (b *Batch) Transform(val any) any {
	b.mx.Lock()
	defer b.mx.Unlock()
	select {
	case b.ch <- val:
		return nil
	default:
		ch := b.ch
		b.ch = make(chan any, b.size)
		close(ch)
		return ch
	}
}
