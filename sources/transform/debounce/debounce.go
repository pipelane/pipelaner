/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package debounce

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/pipeline/source"
)

func init() {
	source.RegisterTransform("debounce", &Debounce{})
}

type Debounce struct {
	interval time.Duration
	val      atomic.Value
	locked   atomic.Bool

	mx    sync.Mutex
	timer *time.Timer
}

func (d *Debounce) Init(cfg transform.Transform) error {
	dConfig, ok := cfg.(transform.Debounce)
	if !ok {
		return fmt.Errorf("invalid debounce config type: %T", cfg)
	}

	d.interval = dConfig.GetInterval().GoDuration()
	d.timer = time.NewTimer(dConfig.GetInterval().GoDuration())
	return nil
}

func (d *Debounce) Transform(val any) any {
	d.storeValue(val)
	lock := d.locked.Load()
	if lock {
		return nil
	}
	d.locked.Store(true)
	<-d.timer.C
	v := d.val.Load()
	d.locked.Store(false)
	d.val = atomic.Value{}
	return v
}

func (d *Debounce) storeValue(val any) {
	d.reset()
	d.val.Store(val)
}

func (d *Debounce) reset() {
	d.mx.Lock()
	defer d.mx.Unlock()
	d.timer.Reset(d.interval)
}
