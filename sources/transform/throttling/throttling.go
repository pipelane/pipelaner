/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package throttling

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/pipeline/source"
	nullableAtomic "github.com/pipelane/pipelaner/sources/shared/atomic"
)

func init() {
	source.RegisterTransform("throttling", &Throttling{})
}

type Throttling struct {
	interval time.Duration
	val      nullableAtomic.Nullable
	locked   atomic.Bool

	mx    sync.Mutex
	timer *time.Timer
}

func (t *Throttling) Init(cfg transform.Transform) error {
	tCfg, ok := cfg.(transform.Throttling)
	if !ok {
		return fmt.Errorf("invalid transform config type: %T", cfg)
	}
	t.interval = tCfg.GetInterval().GoDuration()
	t.timer = time.NewTimer(tCfg.GetInterval().GoDuration())
	return nil
}

func (t *Throttling) Transform(val any) any {
	t.storeValue(val)
	lock := t.locked.Load()
	if lock {
		return nil
	}
	t.locked.Store(true)
	t.reset()
	<-t.timer.C
	v := t.val.Load()
	t.locked.Store(false)
	t.val.Store(nil)
	return v
}

func (t *Throttling) storeValue(val any) {
	t.val.Store(val)
}

func (t *Throttling) reset() {
	t.mx.Lock()
	defer t.mx.Unlock()
	t.timer.Reset(t.interval)
}
