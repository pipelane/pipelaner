package chunks

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
			chunk := NewChunks[int](context.Background(), Config{
				MaxChunkSize: 10,
				BufferSize:   10,
				MaxIdleTime:  time.Second * 10,
			})
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
