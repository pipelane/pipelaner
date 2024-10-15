/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newCfg(
	itemType LaneTypes, //nolint:unparam
	extended map[string]any,
) *BaseLaneConfig {
	c, err := NewBaseConfigWithTypeAndExtended(
		itemType,
		"test_maps_sinks",
		extended,
	)
	if err != nil {
		return nil
	}
	return c
}

func TestSubscriber_Run_Receive(t *testing.T) {
	type args struct {
		iterationsCount int
		threadsCount    int64
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "test iterations 10",
			args: args{
				iterationsCount: 10,
				threadsCount:    1,
			},
			want: []int{
				0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			},
		},
		{
			name: "test iterations 3",
			args: args{
				iterationsCount: 3,
				threadsCount:    100,
			},
			want: []int{
				0, 1, 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inc := 0
			tCount := tt.args.threadsCount
			s := newRunLoop(100, &tCount)
			var res []int
			c, cancel := context.WithCancel(context.Background())
			wg := sync.WaitGroup{}
			mx := sync.Mutex{}
			wg.Add(tt.args.iterationsCount)
			s.methods = methods{
				transform: func(_ *Context, val any) any {
					return val
				},
				sink: func(_ *Context, val any) {
					mx.Lock()
					defer mx.Unlock()
					res = append(res, val.(int))
					wg.Done()
				},
				generator: func(_ *Context, input chan<- any) {
					for {
						if inc < tt.args.iterationsCount {
							input <- inc
							inc++
							continue
						}
						time.Sleep(time.Millisecond * 500)
						cancel()
						return
					}
				},
			}

			ctx := &Context{
				ctx: c,
			}
			s.setContext(ctx)
			s.receive()
			s.start()
			wg.Wait()
			sort.Ints(res)
			assert.Equal(t, res, tt.want)
		})
	}
}

func TestSubscriber_Subscribe(t *testing.T) {
	type args struct {
		iterationsCount int
		threadsCount    int64
		bufferSize      int64
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "test iterations 10",
			args: args{
				iterationsCount: 10,
				threadsCount:    1,
				bufferSize:      100,
			},
			want: []int{
				0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			},
		},
		{
			name: "test iterations 3",
			args: args{
				iterationsCount: 3,
				threadsCount:    100,
				bufferSize:      100,
			},
			want: []int{
				0, 1, 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inc := 0
			tCount := tt.args.threadsCount
			input := newRunLoop(tt.args.bufferSize, &tCount)
			transform := newRunLoop(tt.args.bufferSize, &tCount)
			sink := newRunLoop(tt.args.bufferSize, &tCount)

			var res []int
			c, cancel := context.WithCancel(context.Background())
			wg := sync.WaitGroup{}
			mx := sync.Mutex{}
			wg.Add(tt.args.iterationsCount)
			method := methods{
				transform: func(_ *Context, val any) any {
					return val
				},
				sink: func(_ *Context, val any) {
					mx.Lock()
					defer mx.Unlock()
					res = append(res, val.(int))
					wg.Done()
				},
				generator: func(_ *Context, input chan<- any) {
					for {
						if inc < tt.args.iterationsCount {
							input <- inc
							inc++
							continue
						}
						time.Sleep(time.Millisecond * 500)
						cancel()
						return
					}
				},
			}
			input.setGenerator(method.generator)
			transform.setMap(method.transform)
			sink.setSink(method.sink)

			output := input.createOutput(tt.args.bufferSize)
			transform.setInputChannel(output)
			output = transform.createOutput(tt.args.bufferSize)
			sink.setInputChannel(output)
			ctx := &Context{
				ctx: c,
			}
			input.setContext(ctx)
			transform.setContext(ctx)
			sink.setContext(ctx)
			input.receive()
			input.start()
			transform.start()
			sink.start()
			wg.Wait()
			sort.Ints(res)
			assert.Equal(t, res, tt.want)
		})
	}
}

func TestSubscriber_SubscribeChunks(t *testing.T) {
	type args struct {
		iterationsCount int
		threadsCount    int64
		bufferSize      int64
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test iterations 10",
			args: args{
				iterationsCount: 10,
				threadsCount:    1,
				bufferSize:      100,
			},
			want: []int{
				0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
			},
		},
		{
			name: "test iterations 3",
			args: args{
				iterationsCount: 3,
				threadsCount:    100,
				bufferSize:      100,
			},
			want: []int{
				0, 1, 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inc := 0
			tCount := tt.args.threadsCount
			input := newRunLoop(tt.args.bufferSize, &tCount)
			transform := newRunLoop(tt.args.bufferSize, &tCount)
			sink := newRunLoop(tt.args.bufferSize, &tCount)
			var res []any
			c, cancel := context.WithCancel(context.Background())
			wg := sync.WaitGroup{}
			mx := sync.Mutex{}
			wg.Add(tt.args.iterationsCount)
			method := methods{
				transform: func(ctx *Context, val any) any {
					return val
				},
				sink: func(_ *Context, val any) {
					mx.Lock()
					defer mx.Unlock()
					for v := range val.(chan any) {
						res = append(res, v)
						wg.Done()
					}
				},
				generator: func(_ *Context, input chan<- any) {
					for {
						ch := make(chan int, tt.args.iterationsCount)
						for i := 0; inc < tt.args.iterationsCount; i++ {
							ch <- inc
						}
						input <- ch
						time.Sleep(time.Second * 10)
						cancel()
						return
					}
				},
			}
			input.setGenerator(method.generator)
			transform.setMap(method.transform)
			sink.setSink(method.sink)

			output := input.createOutput(tt.args.bufferSize)
			transform.setInputChannel(output)
			output = transform.createOutput(tt.args.bufferSize)
			sink.setInputChannel(output)
			ctx := &Context{
				ctx: c,
			}
			input.setContext(ctx)
			transform.setContext(ctx)
			sink.setContext(ctx)
			input.receive()
			input.start()
			transform.start()
			sink.start()
			wg.Wait()
			//sort.Ints(res)
			assert.Equal(t, res, tt.want)
		})
	}
}
