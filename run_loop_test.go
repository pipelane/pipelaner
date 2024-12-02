/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pipelane/pipelaner/source/shared/chunker"
	"github.com/stretchr/testify/assert"
)

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
			s := newRunLoop(100, tCount, false, false)
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
			input := newRunLoop(tt.args.bufferSize, tCount, false, false)
			transform := newRunLoop(tt.args.bufferSize, tCount, false, false)
			sink := newRunLoop(tt.args.bufferSize, tCount, false, false)

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
			tCount := tt.args.threadsCount
			input := newRunLoop(tt.args.bufferSize, tCount, false, false)
			transform := newRunLoop(tt.args.bufferSize, tCount, false, false)
			sink := newRunLoop(tt.args.bufferSize, tCount, false, false)
			var res []int
			c, cancel := context.WithCancel(context.Background())
			wg := sync.WaitGroup{}
			mx := sync.Mutex{}
			wg.Add(tt.args.iterationsCount)
			gen := chunker.NewChunks[any](context.Background(), chunker.Config{
				MaxChunkSize: tt.args.iterationsCount,
				BufferSize:   2,
				MaxIdleTime:  time.Second * 10,
			})
			gen.Generator()
			locked := atomic.Bool{}
			method := methods{
				transform: func(_ *Context, val any) any {
					gen.Input() <- val
					if locked.Load() {
						return nil
					}
					locked.Store(true)
					defer locked.Store(false)
					v := <-gen.GetChunks()
					return v
				},
				sink: func(_ *Context, val any) {
					mx.Lock()
					defer mx.Unlock()
					for vch := range val.(chan any) {
						for v := range vch.(chan any) {
							res = append(res, v.(int))
							wg.Done()
						}
					}
				},
				generator: func(_ *Context, input chan<- any) {
					ch := make(chan any, tt.args.iterationsCount)
					for i := 0; i < tt.args.iterationsCount; i++ {
						ch <- i
					}
					input <- ch
					close(ch)
					time.Sleep(time.Second * 10)
					cancel()
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
