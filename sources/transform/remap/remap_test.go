/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package remap

import (
	"testing"

	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExprLanguage_Map(t *testing.T) {
	type args struct {
		val  any
		code string
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test expr maps return nil",
			args: args{
				val: map[string]any{
					"id":       1,
					"name":     "iPhone 12",
					"price":    999,
					"quantity": 1,
				},
				code: "{ \"value_name\": Data.name, \"value_price\": Data.price}",
			},
			want: map[string]any{
				"value_name":  "iPhone 12",
				"value_price": 999,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Remap{}
			err := e.Init(&transform.RemapImpl{
				Code: tt.args.code,
			})
			require.Nil(t, err)
			got := e.Transform(tt.args.val)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestExprLanguage_String(t *testing.T) {
	type args struct {
		val  any
		code string
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test remap string return nil",
			args: args{
				code: "{ \"value_name\": Data.name, \"value_price\": Data.price}",
				val:  "  {\"id\": 1,\"name\": \"iPhone 12\",\"price\": \"999\",\"quantity\": 1}",
			},
			want: map[string]any{
				"value_name":  "iPhone 12",
				"value_price": "999",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Remap{}
			err := e.Init(&transform.RemapImpl{
				Code: tt.args.code,
			})
			require.Nil(t, err)
			got := e.Transform(tt.args.val)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestExprLanguage_StringArray(t *testing.T) {
	type args struct {
		val  any
		code string
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test remap string array",
			args: args{
				code: "{ \"value_name\": Data[0].name, \"value_price\": Data[0].price}",
				val:  "[{\"id\": 1,\"name\": \"iPhone 12\",\"price\": \"999\",\"quantity\": 1}, {\"id\": 2,\"name\": \"iPhone 13\",\"price\": \"999\",\"quantity\": 1}]",
			},
			want: map[string]any{
				"value_name":  "iPhone 12",
				"value_price": "999",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Remap{}
			err := e.Init(&transform.RemapImpl{
				Code: tt.args.code,
			})
			require.Nil(t, err)
			got := e.Transform(tt.args.val)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestExprLanguage_Bytes(t *testing.T) {
	type args struct {
		val  any
		code string
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test remap string return nil",
			args: args{
				code: "{ \"value_name\": Data.name, \"value_price\": Data.price}",
				val:  []byte("{\"id\": 1,\"name\": \"iPhone 12\",\"price\": \"999\",\"quantity\": 1}"),
			},
			want: map[string]any{
				"value_name":  "iPhone 12",
				"value_price": "999",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Remap{}
			err := e.Init(&transform.RemapImpl{
				Code: tt.args.code,
			})
			require.Nil(t, err)
			got := e.Transform(tt.args.val)
			assert.Equal(t, got, tt.want)
		})
	}
}
