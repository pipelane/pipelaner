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

type ThrottlingCfg struct {
	Interval string `pipelane:"interval"`
}

type Throttling struct {
	mx     sync.Mutex
	cfg    *pipelaner.BaseLaneConfig
	val    atomic.Value
	locked atomic.Bool
	timer  *time.Timer
}

func (d *Throttling) Init(cfg *pipelaner.BaseLaneConfig) error {
	d.cfg = cfg
	v := &ThrottlingCfg{}
	err := cfg.ParseExtended(v)
	if err != nil {
		return err
	}
	i, err := d.Interval()
	if err != nil {
		return err
	}
	d.timer = time.NewTimer(i)
	return nil
}

func (d *Throttling) New() pipelaner.Map {
	return &Throttling{
		cfg:   d.cfg,
		timer: d.timer,
		mx:    sync.Mutex{},
	}
}

func (d *Throttling) Map(ctx context.Context, val any) any {
	d.storeValue(val)
	lock := d.locked.Load()
	if lock {
		return nil
	}
	d.locked.Store(true)
	select {
	case <-ctx.Done():
		d.timer.Stop()
		return nil
	case <-d.timer.C:
		v := d.val.Load()
		d.locked.Store(false)
		d.val = atomic.Value{}
		d.reset()
		return v
	}
}

func (d *Throttling) storeValue(val any) {
	d.val.Store(val)
}

func (d *Throttling) reset() {
	d.mx.Lock()
	defer d.mx.Unlock()
	i, _ := d.Interval()
	d.timer.Reset(i)
}

func (d *Throttling) Interval() (time.Duration, error) {
	return time.ParseDuration(d.cfg.Extended.(*ThrottlingCfg).Interval)
}
