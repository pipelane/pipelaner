/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package batch

/*import (
	"context"
	"sync"

	"github.com/pipelane/pipelaner"
)

type Config struct {
	Size int64 `pipelane:"size"`
}

type Batch struct {
	cfg *pipelaner.BaseLaneConfig
	mx  sync.Mutex
	ch  chan any
}

func init() {
	pipelaner.RegisterMap("batch", &Batch{})
}

func (b *Batch) Init(ctx *pipelaner.Context) error {
	b.cfg = ctx.LaneItem().Config()
	v := &Config{}
	err := b.cfg.ParseExtended(v)
	if err != nil {
		return err
	}
	b.ch = make(chan any, b.cfg.Extended.(*Config).Size)
	return nil
}

func (b *Batch) Map(ctx context.Context, val any) any {
	b.mx.Lock()
	defer b.mx.Unlock()
	select {
	case <-ctx.Done():
		return nil
	case b.ch <- val:
		return nil
	default:
		ch := b.ch
		b.ch = make(chan any, b.cfg.Extended.(*Config).Size)
		close(ch)
		return ch
	}
}*/
