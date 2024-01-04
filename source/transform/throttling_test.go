/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package transform

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"pipelaner"
)

func TestThrottling_Map(t *testing.T) {
	type args struct {
		cfg        *pipelaner.BaseLaneConfig
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
				cfg: newCfg(pipelaner.MapType,
					"test_maps",
					map[string]any{
						"interval": "300ms",
					},
				),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Throttling{
				mx:  sync.Mutex{},
				cfg: tt.args.cfg,
				val: atomic.Value{},
			}
			e := d.Init(d.cfg)
			if e != nil {
				t.Error(e)
				return
			}
			maps := d.New()
			var val *int
			for i := 0; i < tt.args.iterations; i++ {
				v := maps.Map(context.Background(), i)
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
		cfg        *pipelaner.BaseLaneConfig
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
				cfg: newCfg(pipelaner.MapType,
					"test_maps",
					map[string]any{
						"interval": "300ms",
					},
				),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Throttling{
				mx:  sync.Mutex{},
				cfg: tt.args.cfg,
				val: atomic.Value{},
			}
			e := d.Init(d.cfg)
			if e != nil {
				t.Error(e)
				return
			}
			maps := d.New()
			wg := sync.WaitGroup{}
			var val *int
			for i := 0; i < tt.args.iterations; i++ {
				wg.Add(1)
				go func(j int) {
					defer wg.Done()
					v := maps.Map(context.Background(), j)
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
