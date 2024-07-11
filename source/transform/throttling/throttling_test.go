/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package throttling

import (
	"context"
	"pipelaner/source/transform/filter"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"pipelaner"
)

func TestThrottling_Map(t *testing.T) {
	type args struct {
		ctx        *pipelaner.Context
		iterations int
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "Test throttling 300 ms",
			args: args{
				iterations: 10,
				ctx: pipelaner.NewContext(
					context.Background(),
					pipelaner.NewLaneItem(filter.newCfg(pipelaner.MapType,
						"test_maps",
						map[string]any{
							"interval": "300ms",
						},
					))),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Throttling{
				mx:  sync.Mutex{},
				cfg: tt.args.ctx.LaneItem().Config(),
				val: atomic.Value{},
			}
			maps := d.New()
			e := maps.Init(tt.args.ctx)
			if e != nil {
				t.Error(e)
				return
			}
			var val *int
			for i := 0; i < tt.args.iterations; i++ {
				v := maps.Map(tt.args.ctx, i)
				if v != nil {
					assert.Equal(t, v, i)
					continue
				}
				assert.Nil(t, val)
			}
		})
	}
}

func TestThrottlingConcurrent_Map(t *testing.T) {
	type args struct {
		ctx        *pipelaner.Context
		iterations int
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "Test concurrent throttling 300 ms",
			args: args{
				iterations: 10,
				ctx: pipelaner.NewContext(
					context.Background(),
					pipelaner.NewLaneItem(filter.newCfg(pipelaner.MapType,
						"test_maps",
						map[string]any{
							"interval": "300ms",
						},
					))),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Throttling{
				mx:  sync.Mutex{},
				cfg: tt.args.ctx.LaneItem().Config(),
				val: atomic.Value{},
			}
			maps := d.New()
			e := maps.Init(tt.args.ctx)
			if e != nil {
				t.Error(e)
				return
			}
			wg := sync.WaitGroup{}
			var val *int
			for i := 0; i < tt.args.iterations; i++ {
				wg.Add(1)
				go func(j int) {
					defer wg.Done()
					v := maps.Map(tt.args.ctx, j)
					if v != nil {
						_v := v.(int)
						val = &_v
					}
				}(i)
			}
			wg.Wait()
			i, _ := d.Interval()
			time.Sleep(i + time.Second)
			assert.NotNil(t, val)
		})
	}
}
