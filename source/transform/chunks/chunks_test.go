package chunks

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestChunks(t *testing.T) {
	type args struct {
		cfg Config
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test chunks 10",
			args: args{
				cfg: Config{
					MaxChunkSize: 10,
					BufferSize:   10,
					MaxIdleTime:  time.Second * 10,
				},
			},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chunk := NewChunks[int](context.Background(), tt.args.cfg)
			chunk.Generator()
			go func() {
				for i := 0; i < tt.args.cfg.BufferSize; i++ {
					chunk.Input() <- i
				}
			}()
			buff := chunk.GetChunks()
			output := <-buff
			var slice []int
			for out := range output {
				slice = append(slice, out)
			}
			assert.Equal(t, tt.want, slice)
		})
	}
}

func TestChunksOfChunks(t *testing.T) {
	type args struct {
		cfg Config
	}
	var tests = []struct {
		name string
		args args
		want []any
	}{
		{
			name: "Test chunks 10",
			args: args{
				cfg: Config{
					MaxChunkSize: 3,
					BufferSize:   3,
					MaxIdleTime:  time.Second * 10,
				},
			},
			want: []any{0, 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			chunk := NewChunks[chan any](ctx, tt.args.cfg)
			chunk.Generator()
			go func() {
				ch := make(chan any, tt.args.cfg.BufferSize)
				for i := 0; i < tt.args.cfg.BufferSize; i++ {
					ch <- i
				}
				chunk.Input() <- ch
				go func() {
					time.Sleep(time.Second * 10)
					close(ch)
				}()
			}()
			buff := chunk.GetChunks()
			output := <-buff
			var slice []any
			for out := range output {
				for v := range out {
					slice = append(slice, v)
				}
			}
			cancel()
			assert.Equal(t, tt.want, slice)
		})
	}
}
