package pipelane

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestSubscriber_Subscribe(t *testing.T) {
	type fields struct {
	}
	type args struct {
		newBufferSize    int64
		subscribersCount int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int
	}{
		{
			name:   "test createOutput 1",
			fields: fields{},
			args: args{
				newBufferSize:    3,
				subscribersCount: 1,
			},
			want: []int{
				0, 1, 2,
			},
		},
		{
			name:   "test createOutput 2",
			fields: fields{},
			args: args{
				newBufferSize:    3,
				subscribersCount: 2,
			},
			want: []int{
				0, 1, 2, 0, 1, 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, cancel := context.WithCancel(context.Background())
			s := newSubscriber(c, tt.args.newBufferSize, nil)
			ch := make(chan any, tt.args.newBufferSize)
			s.setInputChannel(ch)
			gr := errgroup.Group{}
			var res []int
			for i := 0; i < tt.args.subscribersCount; i++ {
				outPut := s.Subscribe(tt.args.newBufferSize)
				gr.Go(func() error {
					for v := range outPut {
						res = append(res, v.(int))
					}
					return nil
				})
			}
			gr.Go(func() error {
				for i := 0; i < int(tt.args.newBufferSize); i++ {
					ch <- i
				}
				return nil
			})
			gr.Go(func() error {
				time.Sleep(time.Second * 1)
				cancel()
				close(ch)
				return nil
			})
			err := gr.Wait()
			assert.Equal(t, res, tt.want)
			assert.Nil(t, err)
		})
	}
}
