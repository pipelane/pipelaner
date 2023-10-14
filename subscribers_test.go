package pipelane

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscriber_Subscribe(t *testing.T) {
	type fields struct {
	}
	type args struct {
		newBufferSize int64
		maxValue      int
		generator     MethodGenerator
		maps          MethodMap
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int
	}{
		{
			name:   "test createOutput and map nil",
			fields: fields{},
			args: args{
				newBufferSize: 3,
				maxValue:      3,
				generator: func() MethodGenerator {
					i := 0
					return func(ctx context.Context) any {
						i++
						return i
					}
				}(),
				maps: nil,
			},
			want: []int{
				1, 2, 3,
			},
		},
		{
			name:   "test createOutput and map",
			fields: fields{},
			args: args{
				newBufferSize: 3,
				maxValue:      4,
				generator: func() MethodGenerator {
					i := 0
					return func(ctx context.Context) any {
						i++
						return i
					}
				}(),
				maps: func(ctx context.Context, val any) any {
					v := val.(int)
					v++
					return v
				},
			},
			want: []int{
				2, 3, 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, cancel := context.WithCancel(context.Background())
			s := newSubscriber(c, tt.args.newBufferSize)
			s.SetGenerator(tt.args.generator)
			var res []int
			s.SetMap(tt.args.maps)
			s.SetSink(func(ctx context.Context, val any) {
				if val.(int) > tt.args.maxValue {
					cancel()
				} else {
					res = append(res, val.(int))
				}
			})
			s.run()
			s.Receive()
			<-c.Done()
			assert.Equal(t, res, tt.want)
		})
	}
}
