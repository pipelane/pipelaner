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
			s := newRunLoop(100, &tt.args.threadsCount)
			var res []int
			c, cancel := context.WithCancel(context.Background())
			wg := sync.WaitGroup{}
			mx := sync.Mutex{}
			wg.Add(tt.args.iterationsCount)
			s.methods = methods{
				transform: func(ctx context.Context, val any) any {
					return val
				},
				sink: func(ctx context.Context, val any) {
					mx.Lock()
					defer mx.Unlock()
					res = append(res, val.(int))
					wg.Done()
				},
				generator: func(ctx context.Context, input chan<- any) {
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

			s.Receive(c)
			s.run(c)
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
			input := newRunLoop(tt.args.bufferSize, &tt.args.threadsCount)
			transform := newRunLoop(tt.args.bufferSize, &tt.args.threadsCount)
			sink := newRunLoop(tt.args.bufferSize, &tt.args.threadsCount)

			var res []int
			c, cancel := context.WithCancel(context.Background())
			wg := sync.WaitGroup{}
			mx := sync.Mutex{}
			wg.Add(tt.args.iterationsCount)
			method := methods{
				transform: func(ctx context.Context, val any) any {
					return val
				},
				sink: func(ctx context.Context, val any) {
					mx.Lock()
					defer mx.Unlock()
					res = append(res, val.(int))
					wg.Done()
				},
				generator: func(ctx context.Context, input chan<- any) {
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
			input.SetGenerator(method.generator)
			transform.SetMap(method.transform)
			sink.SetSink(method.sink)

			output := input.createOutput(tt.args.bufferSize)
			transform.setInputChannel(output)
			output = transform.createOutput(tt.args.bufferSize)
			sink.setInputChannel(output)

			input.Receive(c)
			input.run(c)
			transform.run(c)
			sink.run(c)
			wg.Wait()
			sort.Ints(res)
			assert.Equal(t, res, tt.want)
		})
	}
}
