/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package transform

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"pipelaner"
)

type DebounceCfg struct {
	Interval string `pipelane:"interval"`
}

type Debounce struct {
	mx     sync.Mutex
	cfg    *pipelaner.BaseLaneConfig
	val    atomic.Value
	locked atomic.Bool
	timer  *time.Timer
}

func (d *Debounce) Init(cfg *pipelaner.BaseLaneConfig) error {
	d.cfg = cfg
	v := &DebounceCfg{}
	err := cfg.ParseExtended(v)
	if err != nil {
		return err
	}
	interval, err := d.Interval()
	if err != nil {
		return err
	}
	d.timer = time.NewTimer(interval)
	return nil
}

func (d *Debounce) New() pipelaner.Map {
	return &Debounce{
		cfg:   d.cfg,
		timer: d.timer,
		mx:    sync.Mutex{},
	}
}

func (d *Debounce) Map(ctx context.Context, val any) any {
	d.storeValue(val)
	lock := d.locked.Load()
	if lock {
		return nil
	}
	d.locked.Store(true)
	select {
	case <-ctx.Done():
		return nil
	case <-d.timer.C:
		v := d.val.Load()
		d.locked.Store(false)
		d.val = atomic.Value{}
		d.reset()
		return v
	}
}

func (d *Debounce) storeValue(val any) {
	d.reset()
	d.val.Store(val)
}

func (d *Debounce) reset() {
	d.mx.Lock()
	defer d.mx.Unlock()
	i, _ := d.Interval()
	d.timer.Reset(i)
}

func (d *Debounce) Interval() (time.Duration, error) {
	return time.ParseDuration(d.cfg.Extended.(*DebounceCfg).Interval)
}