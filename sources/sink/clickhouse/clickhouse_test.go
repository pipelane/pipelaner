/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package clickhouse

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ClickHouse/ch-go/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuildProtoInput(t *testing.T) {
	timeValue := time.Now()

	type fields struct {
	}
	type args struct {
		values map[string]any
	}

	type want struct {
		mapColumns map[string]*column
		input      proto.Input
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    want
	}{
		{
			name:   "build success",
			fields: fields{},
			args: args{
				values: map[string]any{
					"string":       "123",
					"string_array": []string{"123", "456"},
					"float":        1.0,
					"time":         timeValue,
				},
			},
			wantErr: false,
			want: want{
				mapColumns: map[string]*column{
					"string":       {str: new(proto.ColStr)},
					"string_array": {strArr: new(proto.ColArr[string])},
					"float":        {flt: new(proto.ColFloat64)},
					"time":         {timestamp: new(proto.ColDateTime64)},
				},
				input: proto.Input{
					proto.InputColumn{Name: "string", Data: new(proto.ColStr)},
					proto.InputColumn{Name: "string_array", Data: new(proto.ColArr[string])},
					proto.InputColumn{Name: "float", Data: new(proto.ColFloat64)},
					proto.InputColumn{Name: "time", Data: new(proto.ColDateTime64)},
				},
			},
		},
		{
			name:   "build slice boolean",
			fields: fields{},
			args: args{
				values: map[string]any{
					"bool_array":   []bool{true},
					"string_array": []string{"123", "456"},
					"float":        1.0,
					"time":         timeValue,
				},
			},
			want: want{
				mapColumns: map[string]*column{
					"bool_array":   {boolArr: new(proto.ColArr[bool])},
					"string_array": {strArr: new(proto.ColArr[string])},
					"float":        {flt: new(proto.ColFloat64)},
					"time":         {timestamp: new(proto.ColDateTime64)},
				},
				input: proto.Input{
					proto.InputColumn{Name: "bool_array", Data: new(proto.ColArr[bool])},
					proto.InputColumn{Name: "string_array", Data: new(proto.ColArr[string])},
					proto.InputColumn{Name: "float", Data: new(proto.ColFloat64)},
					proto.InputColumn{Name: "time", Data: new(proto.ColDateTime64)},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := Clickhouse{
				client: &Client{},
			}

			mapColumns, input, err := cl.buildProtoInput(tt.args.values)
			if tt.wantErr {
				assert.Nil(t, mapColumns)
				assert.Nil(t, input)
				assert.Error(t, err)
				return
			}
			EqualInput(t, tt.want.input, input)
			EqualNotNilKeys(t, tt.want.mapColumns, mapColumns)

			assert.Nil(t, err)
		})
	}
}

func EqualInput(t *testing.T, a, b proto.Input) bool {
	t.Helper()
	if len(a) != len(b) {
		assert.Fail(t, "length input not equal")
	}

loop:
	for _, v := range a {
		for _, vv := range b {
			if vv.Name == v.Name {
				if v.Data == nil || vv.Data == nil || reflect.ValueOf(v.Data).Type() != reflect.ValueOf(vv.Data).Type() {
					assert.Fail(t, fmt.Sprintf("Not equal: \n"+
						"expected: name=%s%+v"+
						"actual  :name=%s %+v", v.Name, v.Data, vv.Name, vv.Data))
				}

				continue loop
			}
		}
		assert.Fail(t, "input not equal")
	}

	return true
}

func EqualNotNilKeys(t *testing.T, a, b map[string]*column) bool {
	t.Helper()
	if len(a) != len(b) {
		assert.Fail(t, "length transform not equal")
	}

	for k, v := range a {
		if val, ok := b[k]; !ok || v == nil || val == nil || reflect.ValueOf(v).Type() != reflect.ValueOf(val).Type() {
			assert.Fail(t, fmt.Sprintf("Not equal: \n"+
				"expected: %+v"+
				"actual  : %+v", v, val))
		}
	}

	return true
}
