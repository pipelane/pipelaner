/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package debounce

import (
	"sync"
	"testing"

	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/stretchr/testify/assert"
)

func newCfg(
	t *testing.T,
	duration float64,
	durationUnit string,
) transform.Transform {
	t.Helper()
	uni, err := pkl.ToDurationUnit(durationUnit)
	if err != nil {
		t.Fatal(err)
	}
	return &transform.DebounceImpl{
		Interval: pkl.Duration{
			Unit:  uni,
			Value: duration,
		},
	}
}

func TestDebounce_Map(t *testing.T) {
	type args struct {
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
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maps := &Debounce{
				mx: sync.Mutex{},
			}
			e := maps.Init(newCfg(t, 300, "ms"))
			if e != nil {
				t.Error(e)
				return
			}
			var val *int
			for i := 0; i < tt.args.iterations; i++ {
				v := maps.Transform(i)
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
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maps := &Debounce{
				mx: sync.Mutex{},
			}
			e := maps.Init(newCfg(t, 300, "ms"))
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
					v := maps.Transform(j)
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
			assert.NotNil(t, val)
		})
	}
}
