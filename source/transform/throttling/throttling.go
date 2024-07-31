/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package throttling

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/pipelane/pipelaner"
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

func (d *Throttling) Init(ctx *pipelaner.Context) error {
	d.cfg = ctx.LaneItem().Config()
	v := &ThrottlingCfg{}
	err := d.cfg.ParseExtended(v)
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

func init() {
	pipelaner.RegisterMap("throttling", &Throttling{})
}

func (d *Throttling) Map(ctx *pipelaner.Context, val any) any {
	d.storeValue(val)
	lock := d.locked.Load()
	if lock {
		return nil
	}
	d.locked.Store(true)
	d.reset()
	select {
	case <-ctx.Context().Done():
		d.timer.Stop()
		return nil
	case <-d.timer.C:
		v := d.val.Load()
		d.locked.Store(false)
		d.val = atomic.Value{}
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
