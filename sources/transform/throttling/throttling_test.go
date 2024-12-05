/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package throttling

import (
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/stretchr/testify/assert"
)

func newConfig(
	t *testing.T,
	duration time.Duration,
) transform.Transform {
	t.Helper()
	durationStr := duration.String()

	re := regexp.MustCompile(`^(\d+)([a-z]+)$`)
	matches := re.FindStringSubmatch(durationStr)
	if len(matches) == 3 {
		value, err := strconv.ParseFloat(matches[1], 64)
		assert.NoError(t, err)
		unit, err := pkl.ToDurationUnit(matches[2])
		assert.NoError(t, err)
		return &transform.ThrottlingImpl{
			Interval: &pkl.Duration{
				Value: value,
				Unit:  unit,
			},
		}
	}
	t.Fatal("invalid input string")
	return nil
}

func TestThrottling_Map(t *testing.T) {
	type args struct {
		iterations int
		duration   time.Duration
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
				duration:   300 * time.Millisecond,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maps := &Throttling{
				mx:  sync.Mutex{},
				val: atomic.Value{},
			}
			e := maps.Init(newConfig(t, tt.args.duration))
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

func TestThrottlingConcurrent_Map(t *testing.T) {
	type args struct {
		iterations int
		duration   time.Duration
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
				duration:   300 * time.Millisecond,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maps := &Throttling{
				mx:  sync.Mutex{},
				val: atomic.Value{},
			}
			e := maps.Init(newConfig(t, tt.args.duration))
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
			time.Sleep(tt.args.duration)
			assert.NotNil(t, val)
		})
	}
}
