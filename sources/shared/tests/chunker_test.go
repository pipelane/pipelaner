package tests

import (
	"testing"
	"time"

	"github.com/pipelane/pipelaner/sources/shared/chunker"
	"github.com/stretchr/testify/assert"
)

func TestChunks(t *testing.T) {
	type args struct {
		cfg chunker.Config
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		{
			name: "Test chunks 10",
			args: args{
				cfg: chunker.Config{
					MaxChunkSize: 10,
					BufferSize:   10,
					MaxIdleTime:  time.Second * 10,
				},
			},
			want: []any{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunk := chunker.NewChunks(tt.args.cfg)
			chunk.Generator()
			go func() {
				for i := 0; i < tt.args.cfg.BufferSize; i++ {
					chunk.SetValue(i)
				}
				chunk.Stop()
			}()
			output := chunk.Chunk()
			var slice []any
			for out := range output {
				slice = append(slice, out)
			}
			assert.Equal(t, tt.want, slice)
		})
	}
}

func TestChunksOfChunks(t *testing.T) {
	type args struct {
		cfg chunker.Config
	}
	type testStructs struct {
		val []any
	}
	var tests = []struct {
		name string
		args args
		want any
	}{
		{
			name: "Test chunks of structs",
			args: args{
				cfg: chunker.Config{
					MaxChunkSize: 3,
					BufferSize:   3,
					MaxIdleTime:  time.Second * 10,
				},
			},
			want: testStructs{
				val: []any{0, 1, 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunk := chunker.NewChunks(tt.args.cfg)
			chunk.Generator()
			go func() {
				ch := make(chan any, tt.args.cfg.BufferSize)
				var slice []any
				tStruct := testStructs{}
				for i := 0; i < tt.args.cfg.BufferSize; i++ {
					slice = append(slice, i)
				}
				tStruct.val = slice
				ch <- tStruct
				chunk.SetValue(ch)
				close(ch)
				chunk.Stop()
			}()
			output := chunk.Chunk()
			var got testStructs
			var ok bool
			for out := range output {
				for v := range out.(chan any) {
					got, ok = v.(testStructs)
					assert.True(t, ok)
				}
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
