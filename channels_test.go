/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mergeInputs(t *testing.T) {
	type args[T any] struct {
		ctx context.Context
		chs []chan T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want int
	}
	tests := []testCase[int]{
		{
			name: "test 1:2 len = 3",
			args: args[int]{
				ctx: context.Background(),
				chs: []chan int{
					make(chan int, 1),
					make(chan int, 2),
				},
			},
			want: 3,
		},
		{
			name: "test 1:1 len = 2",
			args: args[int]{
				ctx: context.Background(),
				chs: []chan int{
					make(chan int, 1),
					make(chan int, 1),
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cap(mergeInputs(tt.args.ctx, tt.args.chs...)))
		})
	}
}
