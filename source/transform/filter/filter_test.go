/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package remap

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pipelane/pipelaner"
)

func newCfg(
	itemType pipelaner.LaneTypes, //nolint:unparam
	extended map[string]any,
) *pipelaner.BaseLaneConfig {
	c, err := pipelaner.NewBaseConfigWithTypeAndExtended(
		itemType,
		"test_maps_sinks",
		extended,
	)
	if err != nil {
		return nil
	}
	return c
}

func TestExprLanguage_Map(t *testing.T) {
	type args struct {
		val any
		ctx *pipelaner.Context
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test filtering maps return nil",
			args: args{
				ctx: pipelaner.NewContext(context.Background(),
					pipelaner.NewLaneItem(newCfg(pipelaner.MapType,
						map[string]any{
							"code": "Data.count > 5",
						}), true),
				),
				val: map[string]any{
					"count": 1,
				},
			},
			want: nil,
		},
		{
			name: "test filtering maps return 10",
			args: args{
				ctx: pipelaner.NewContext(context.Background(),
					pipelaner.NewLaneItem(newCfg(pipelaner.MapType,
						map[string]any{
							"code": "Data.count > 5",
						}), true),
				),
				val: map[string]any{
					"count": 10,
				},
			},
			want: map[string]any{
				"count": 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.Nop()
			e := &Filter{
				logger: &logger,
			}
			err := e.Init(tt.args.ctx)
			require.Nil(t, err)
			got := e.Map(tt.args.ctx, tt.args.val)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestExprLanguage_String(t *testing.T) {
	type args struct {
		val any
		ctx *pipelaner.Context
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test filtering string return nil",
			args: args{
				ctx: pipelaner.NewContext(context.Background(),
					pipelaner.NewLaneItem(newCfg(pipelaner.MapType,
						map[string]any{
							"code": "Data.count > 5",
						}), true),
				),
				val: "{\"count\":1}",
			},
			want: nil,
		},
		{
			name: "test filtering string return 10",
			args: args{
				ctx: pipelaner.NewContext(context.Background(),
					pipelaner.NewLaneItem(newCfg(pipelaner.MapType,
						map[string]any{
							"code": "Data.count > 5",
						}), true),
				),
				val: "{\"count\":10}",
			},
			want: "{\"count\":10}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.Nop()
			e := &Filter{
				logger: &logger,
			}
			err := e.Init(tt.args.ctx)
			require.Nil(t, err)
			got := e.Map(tt.args.ctx, tt.args.val)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestExprLanguage_Bytes(t *testing.T) {
	type args struct {
		val any
		ctx *pipelaner.Context
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test filtering string return nil",
			args: args{
				ctx: pipelaner.NewContext(context.Background(),
					pipelaner.NewLaneItem(newCfg(pipelaner.MapType,
						map[string]any{
							"code": "Data.count > 5",
						}), true),
				),
				val: []byte("{\"count\":1}"),
			},
			want: nil,
		},
		{
			name: "test filtering string return 10",
			args: args{
				ctx: pipelaner.NewContext(context.Background(),
					pipelaner.NewLaneItem(newCfg(pipelaner.MapType,
						map[string]any{
							"code": "Data.count > 5",
						}), true),
				),
				val: []byte("{\"count\":10}"),
			},
			want: []byte("{\"count\":10}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.Nop()
			e := &Filter{
				logger: &logger,
			}
			err := e.Init(tt.args.ctx)
			require.Nil(t, err)
			got := e.Map(tt.args.ctx, tt.args.val)
			assert.Equal(t, got, tt.want)
		})
	}
}
