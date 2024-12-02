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
			name: "test expr maps return nil",
			args: args{
				ctx: pipelaner.NewContext(context.Background(),
					pipelaner.NewLaneItem(newCfg(pipelaner.MapType,
						map[string]any{
							"code": "{ \"value_name\": Data.name, \"value_price\": Data.price}",
						}), true),
				),
				val: map[string]any{
					"id":       1,
					"name":     "iPhone 12",
					"price":    999,
					"quantity": 1,
				},
			},
			want: map[string]any{
				"value_name":  "iPhone 12",
				"value_price": 999,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.Nop()
			e := &Remap{
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
			name: "test remap string return nil",
			args: args{
				ctx: pipelaner.NewContext(context.Background(),
					pipelaner.NewLaneItem(newCfg(pipelaner.MapType,
						map[string]any{
							"code": "{ \"value_name\": Data.name, \"value_price\": Data.price}",
						}), true),
				),
				val: "  {\"id\": 1,\"name\": \"iPhone 12\",\"price\": \"999\",\"quantity\": 1}",
			},
			want: map[string]any{
				"value_name":  "iPhone 12",
				"value_price": "999",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.Nop()
			e := &Remap{
				logger: &logger,
			}
			err := e.Init(tt.args.ctx)
			require.Nil(t, err)
			got := e.Map(tt.args.ctx, tt.args.val)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestExprLanguage_StringArray(t *testing.T) {
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
			name: "test remap string array",
			args: args{
				ctx: pipelaner.NewContext(context.Background(),
					pipelaner.NewLaneItem(newCfg(pipelaner.MapType,
						map[string]any{
							"code": "{ \"value_name\": Data[0].name, \"value_price\": Data[0].price}",
						}), true),
				),
				val: "[{\"id\": 1,\"name\": \"iPhone 12\",\"price\": \"999\",\"quantity\": 1}, {\"id\": 2,\"name\": \"iPhone 13\",\"price\": \"999\",\"quantity\": 1}]",
			},
			want: map[string]any{
				"value_name":  "iPhone 12",
				"value_price": "999",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.Nop()
			e := &Remap{
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
			name: "test remap string return nil",
			args: args{
				ctx: pipelaner.NewContext(context.Background(),
					pipelaner.NewLaneItem(newCfg(pipelaner.MapType,
						map[string]any{
							"code": "{ \"value_name\": Data.name, \"value_price\": Data.price}",
						}), true),
				),
				val: []byte("{\"id\": 1,\"name\": \"iPhone 12\",\"price\": \"999\",\"quantity\": 1}"),
			},
			want: map[string]any{
				"value_name":  "iPhone 12",
				"value_price": "999",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.Nop()
			e := &Remap{
				logger: &logger,
			}
			err := e.Init(tt.args.ctx)
			require.Nil(t, err)
			got := e.Map(tt.args.ctx, tt.args.val)
			assert.Equal(t, got, tt.want)
		})
	}
}
