/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package debounce

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pipelane/pipelaner"
)

func newCfg(
	itemType pipelaner.LaneTypes,
	name string,
	extended map[string]any,
) *pipelaner.BaseLaneConfig {
	c, err := pipelaner.NewBaseConfigWithTypeAndExtended(
		itemType,
		name,
		extended,
	)
	if err != nil {
		return nil
	}
	return c
}

func TestDebounce_Map(t *testing.T) {
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
			name: "Test debounce 300 ms",
			args: args{
				iterations: 10,
				ctx: pipelaner.NewContext(
					context.Background(),
					pipelaner.NewLaneItem(newCfg(pipelaner.MapType, "test_maps", map[string]any{
						"interval": "300ms",
					}), true)),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maps := &Debounce{
				mx:  sync.Mutex{},
				cfg: tt.args.ctx.LaneItem().Config(),
				val: atomic.Value{},
			}
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

func TestDebounceConcurrent_Map(t *testing.T) {
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
			name: "Test concurrent debounce 300 ms",
			args: args{
				iterations: 10,
				ctx: pipelaner.NewContext(
					context.Background(),
					pipelaner.NewLaneItem(newCfg(pipelaner.MapType, "test_maps", map[string]any{
						"interval": "300ms",
					}), true)),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maps := &Debounce{
				mx:  sync.Mutex{},
				cfg: tt.args.ctx.LaneItem().Config(),
				val: atomic.Value{},
			}
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
						_v, ok := v.(int)
						if !ok {
							assert.Fail(t, "value is not int")
						}
						val = &_v
					}
				}(i)
			}
			wg.Wait()
			i, err := maps.Interval()
			assert.NotNil(t, err)
			time.Sleep(i + time.Second)
			assert.NotNil(t, val)
		})
	}
}
