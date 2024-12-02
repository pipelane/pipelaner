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
			name: "test filtering maps return nil",
			args: args{
				code: "Data.count > 5",
				val: map[string]any{
					"count": 1,
				},
			},
			want: nil,
		},
		{
			name: "test filtering maps return 10",
			args: args{
				code: "Data.count > 5",
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
			e := &Filter{}
			err := e.Init(&transform.FilterImpl{
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
			name: "test filtering string return nil",
			args: args{
				code: "Data.count > 5",
				val:  "{\"count\":1}",
			},
			want: nil,
		},
		{
			name: "test filtering string return 10",
			args: args{
				code: "Data.count > 5",
				val:  "{\"count\":10}",
			},
			want: "{\"count\":10}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Filter{}
			err := e.Init(&transform.FilterImpl{
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
			name: "test filtering string return nil",
			args: args{
				code: "Data.count > 5",
				val:  []byte("{\"count\":1}"),
			},
			want: nil,
		},
		{
			name: "test filtering string return 10",
			args: args{
				code: "Data.count > 5",
				val:  []byte("{\"count\":10}"),
			},
			want: []byte("{\"count\":10}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Filter{}
			err := e.Init(&transform.FilterImpl{
				Code: tt.args.code,
			})
			require.Nil(t, err)
			got := e.Transform(tt.args.val)
			assert.Equal(t, got, tt.want)
		})
	}
}
