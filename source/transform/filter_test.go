/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package transform

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pipelaner"
)

func newCfg(
	itemType pipelaner.LaneTypes,
	name string,
	extended map[string]any,
) *pipelaner.BaseLaneConfig {
	c, _ := pipelaner.NewBaseConfigWithTypeAndExtended(
		itemType,
		name,
		extended,
	)
	return c
}

func TestExprLanguage_Map(t *testing.T) {
	type args struct {
		val any
		cfg *pipelaner.BaseLaneConfig
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test filtering maps return nil",
			args: args{
				cfg: newCfg(pipelaner.SinkType,
					"test_maps_sinks",
					map[string]any{
						"code": "Data.count > 5",
					},
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
				cfg: newCfg(pipelaner.MapType,
					"test_maps",
					map[string]any{
						"code": "Data.count > 5",
					},
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
			e := &Filter{
				logger: zerolog.Nop(),
			}
			err := e.Init(tt.args.cfg)
			require.Nil(t, err)
			got := e.Map(context.Background(), tt.args.val)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestExprLanguage_String(t *testing.T) {
	type args struct {
		val any
		cfg *pipelaner.BaseLaneConfig
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test filtering string return nil",
			args: args{
				cfg: newCfg(pipelaner.SinkType,
					"test_maps_sinks",
					map[string]any{
						"code": "Data.count > 5.0",
					},
				),
				val: "{\"count\":1}",
			},
			want: nil,
		},
		{
			name: "test filtering string return 10",
			args: args{
				cfg: newCfg(pipelaner.MapType,
					"test_maps",
					map[string]any{
						"code": "Data.count > 5.0",
					},
				),
				val: "{\"count\":10}",
			},
			want: "{\"count\":10}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Filter{
				logger: zerolog.Nop(),
			}
			err := e.Init(tt.args.cfg)
			require.Nil(t, err)
			got := e.Map(context.Background(), tt.args.val)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestExprLanguage_Bytes(t *testing.T) {
	type args struct {
		val any
		cfg *pipelaner.BaseLaneConfig
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test filtering string return nil",
			args: args{
				cfg: newCfg(pipelaner.SinkType,
					"test_maps_sinks",
					map[string]any{
						"code": "Data.count > 5.0",
					},
				),
				val: []byte("{\"count\":1}"),
			},
			want: nil,
		},
		{
			name: "test filtering string return 10",
			args: args{
				cfg: newCfg(pipelaner.MapType,
					"test_maps",
					map[string]any{
						"code": "Data.count > 5.0",
					},
				),
				val: []byte("{\"count\":10}"),
			},
			want: []byte("{\"count\":10}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Filter{
				logger: zerolog.Nop(),
			}
			err := e.Init(tt.args.cfg)
			require.Nil(t, err)
			got := e.Map(context.Background(), tt.args.val)
			assert.Equal(t, got, tt.want)
		})
	}
}
