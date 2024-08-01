/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func stringPtr(s string) *string {
	return &s
}

func Test_newConfig(t *testing.T) {
	type args struct {
		tomlString string
	}
	tests := []struct {
		name      string
		args      args
		want      *config
		wantError bool
	}{
		{
			name: "test inputs",
			args: args{
				tomlString: `
[input.input1]
buffer = 1
source_name = "int"
	`,
			},
			want: &config{
				Input: map[string]any{
					"input1": map[string]any{
						"buffer":      int64(1),
						"source_name": "int",
					},
				},
				Map:  nil,
				Sink: nil,
			},
			wantError: false,
		},
		{
			name: "test maps",
			args: args{
				tomlString: `
[map.map2]
buffer = 1
source_name = "int"
	`,
			},
			want: &config{
				Input: nil,
				Map: map[string]any{
					"map2": map[string]any{
						"buffer":      int64(1),
						"source_name": "int",
					},
				},
				Sink: nil,
			},
			wantError: false,
		},
		{
			name: "test sinks",
			args: args{
				tomlString: `
[sink.sink3]
buffer = 1
source_name = "int"
	`,
			},
			want: &config{
				Input: nil,
				Map:   nil,
				Sink: map[string]any{
					"sink3": map[string]any{
						"buffer":      int64(1),
						"source_name": "int",
					},
				},
			},
			wantError: false,
		},
		{
			name: "test inputs sinks map",
			args: args{
				tomlString: `
[input.input1]
buffer = 10
source_name = "input_int"
[map.map2]
buffer = 20
source_name = "map_int"
[sink.sink3]
buffer = 30
source_name = "sink_int"
	`,
			},
			want: &config{
				Input: map[string]any{
					"input1": map[string]any{
						"buffer":      int64(10),
						"source_name": "input_int",
					},
				},
				Map: map[string]any{
					"map2": map[string]any{
						"buffer":      int64(20),
						"source_name": "map_int",
					},
				},
				Sink: map[string]any{
					"sink3": map[string]any{
						"buffer":      int64(30),
						"source_name": "sink_int",
					},
				},
			},
			wantError: false,
		},
		{
			name: "test error",
			args: args{
				tomlString: `
[input.input1
buffer = 10
source_name = "input_int"
[map.map2]
buffer = 20
source_name = "map_int"
[sink.sink3]
buffer = 30
source_name = "sink_int"
	`,
			},
			want:      nil,
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tom, err := decodeTomlString(tt.args.tomlString)
			if tt.wantError && err != nil {
				assert.Error(t, err)
				return
			}
			got, err := newConfig(tom)
			if tt.wantError && err != nil {
				assert.Error(t, err)
				return
			}
			assert.Equalf(t, tt.want, got, "newConfig(%v)", tt.args.tomlString)
		})
	}
}

func Test_NewBaseConfigWithTypeAndExtended(t *testing.T) {
	type args struct {
		itemType LaneTypes
		name     string
		extended map[string]any
	}
	tests := []struct {
		name string
		args args
		want *BaseLaneConfig
	}{
		{
			name: "test input type",
			args: args{
				itemType: InputType,
				name:     "input1",
				extended: map[string]any{
					"buffer":      int64(1),
					"source_name": "int",
					"host":        "0.0.0.0",
					"port":        "8080",
				},
			},
			want: &BaseLaneConfig{
				BufferSize: 1,
				Threads:    nil,
				SourceName: "int",
				Inputs:     nil,
				Internal: Internal{
					Name:     "input1",
					LaneType: InputType,
					Extended: nil,
					_extended: map[string]any{
						"buffer":      int64(1),
						"source_name": "int",
						"host":        "0.0.0.0",
						"port":        "8080",
					},
				},
			},
		},
		{
			name: "test input type with buffer skip buffer",
			args: args{
				itemType: InputType,
				name:     "input1",
				extended: map[string]any{
					"source_name": "int",
					"host":        "0.0.0.0",
					"port":        "8080",
				},
			},
			want: &BaseLaneConfig{
				BufferSize: 1,
				Threads:    nil,
				SourceName: "int",
				Inputs:     nil,
				Internal: Internal{
					Name:     "input1",
					LaneType: InputType,
					Extended: nil,
					_extended: map[string]any{
						"source_name": "int",
						"host":        "0.0.0.0",
						"port":        "8080",
					},
				},
			},
		},
		{
			name: "test input type with buffer zero",
			args: args{

				itemType: InputType,
				name:     "input1",
				extended: map[string]any{
					"buffer":      int64(0),
					"source_name": "int",
					"host":        "0.0.0.0",
					"port":        "8080",
				},
			},
			want: &BaseLaneConfig{
				BufferSize: 1,
				Threads:    nil,
				SourceName: "int",
				Inputs:     nil,
				Internal: Internal{
					Name:     "input1",
					LaneType: InputType,
					Extended: nil,
					_extended: map[string]any{
						"buffer":      int64(0),
						"source_name": "int",
						"host":        "0.0.0.0",
						"port":        "8080",
					},
				},
			},
		},
		{
			name: "test sink type",
			args: args{
				itemType: SinkType,
				name:     "sink1",
				extended: map[string]any{
					"inputs":      []string{"input_name"},
					"source_name": "int",
					"host":        "0.0.0.0",
					"port":        "8080",
				},
			},
			want: &BaseLaneConfig{
				BufferSize: 1,
				Threads:    nil,
				SourceName: "int",
				Inputs:     []string{"input_name"},
				Internal: Internal{
					Name:     "sink1",
					LaneType: SinkType,
					Extended: nil,
					_extended: map[string]any{
						"source_name": "int",
						"host":        "0.0.0.0",
						"port":        "8080",
						"inputs":      []string{"input_name"},
					},
				},
			},
		},
		{
			name: "test map type",
			args: args{
				itemType: MapType,
				name:     "map1",
				extended: map[string]any{
					"inputs":      []string{"input_name"},
					"source_name": "int",
					"host":        "0.0.0.0",
					"port":        "8080",
				},
			},
			want: &BaseLaneConfig{
				BufferSize: 1,
				Threads:    nil,
				SourceName: "int",
				Inputs:     []string{"input_name"},
				Internal: Internal{
					Name:     "map1",
					LaneType: MapType,
					Extended: nil,
					_extended: map[string]any{
						"source_name": "int",
						"host":        "0.0.0.0",
						"port":        "8080",
						"inputs":      []string{"input_name"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBaseConfigWithTypeAndExtended(tt.args.itemType, tt.args.name, tt.args.extended)
			if err != nil {
				assert.NotNil(t, err)
				return
			}
			assert.Equalf(t, tt.want, got, "NewBaseConfigWithTypeAndExtended(%v, %v, %v)", tt.args.itemType, tt.args.name, tt.args.extended)
		})
	}
}

func Test_NewBaseConfigWithTypeAndExtendedFromToml(t *testing.T) {
	type args struct {
		tomlString string
		itemType   LaneTypes
		name       string
	}
	tests := []struct {
		name string
		args args
		want *BaseLaneConfig
	}{
		{
			name: "test input type",
			args: args{
				tomlString: `
[input.input1]
buffer = 1
source_name = "int"
host = "0.0.0.0"
port = "8080"
`,
				itemType: InputType,
				name:     "input1",
			},
			want: &BaseLaneConfig{
				BufferSize: 1,
				Threads:    nil,
				SourceName: "int",
				Inputs:     nil,
				Internal: Internal{
					Name:     "input1",
					LaneType: InputType,
					Extended: nil,
					_extended: map[string]any{
						"buffer":      int64(1),
						"source_name": "int",
						"host":        "0.0.0.0",
						"port":        "8080",
					},
				},
			},
		},
		{
			name: "test input type with buffer skip buffer",
			args: args{
				tomlString: `
[input.input1]
source_name = "int"
host = "0.0.0.0"
port = "8080"
`,
				itemType: InputType,
				name:     "input1",
			},
			want: &BaseLaneConfig{
				BufferSize: 1,
				Threads:    nil,
				SourceName: "int",
				Inputs:     nil,
				Internal: Internal{
					Name:     "input1",
					LaneType: InputType,
					Extended: nil,
					_extended: map[string]any{
						"source_name": "int",
						"host":        "0.0.0.0",
						"port":        "8080",
					},
				},
			},
		},
		{
			name: "test input type with buffer zero",
			args: args{
				tomlString: `
[input.input1]
buffer = 0
source_name = "int"
host = "0.0.0.0"
port = "8080"
`,
				itemType: InputType,
				name:     "input1",
			},
			want: &BaseLaneConfig{
				BufferSize: 1,
				Threads:    nil,
				SourceName: "int",
				Inputs:     nil,
				Internal: Internal{
					Name:     "input1",
					LaneType: InputType,
					Extended: nil,
					_extended: map[string]any{
						"buffer":      int64(0),
						"source_name": "int",
						"host":        "0.0.0.0",
						"port":        "8080",
					},
				},
			},
		},
		{
			name: "test sink type",
			args: args{
				tomlString: `
[sink.sink1]
source_name = "int"
host = "0.0.0.0"
port = "8080"
inputs = ["input_name"]
`,
				itemType: SinkType,
				name:     "sink1",
			},
			want: &BaseLaneConfig{
				BufferSize: 1,
				Threads:    nil,
				SourceName: "int",
				Inputs:     []string{"input_name"},
				Internal: Internal{
					Name:     "sink1",
					LaneType: SinkType,
					Extended: nil,
					_extended: map[string]any{
						"source_name": "int",
						"host":        "0.0.0.0",
						"port":        "8080",
						"inputs":      []any{"input_name"},
					},
				},
			},
		},
		{
			name: "test map type",
			args: args{
				tomlString: `
[map.map1]
source_name = "int"
host = "0.0.0.0"
port = "8080"
inputs = ["input_name"]
`,
				itemType: MapType,
				name:     "map1",
			},
			want: &BaseLaneConfig{
				BufferSize: 1,
				Threads:    nil,
				SourceName: "int",
				Inputs:     []string{"input_name"},
				Internal: Internal{
					Name:     "map1",
					LaneType: MapType,
					Extended: nil,
					_extended: map[string]any{
						"source_name": "int",
						"host":        "0.0.0.0",
						"port":        "8080",
						"inputs":      []any{"input_name"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tom, err := decodeTomlString(tt.args.tomlString)
			if err != nil {
				assert.Error(t, err)
				return
			}
			data, ok := tom[string(tt.args.itemType)].(map[string]any)
			if !ok {
				assert.Error(t, errors.New("not a map[string]any"))
				return
			}
			extended, ok := data[tt.args.name].(map[string]any)
			if !ok {
				assert.Error(t, errors.New("not a map[string]any"))
				return
			}
			got, err := NewBaseConfigWithTypeAndExtended(tt.args.itemType, tt.args.name, extended)
			if err != nil {
				assert.NotNil(t, err)
				return
			}
			assert.Equalf(t, tt.want, got, "NewBaseConfigWithTypeAndExtended(%v, %v, %v)", tt.args.itemType, tt.args.name, extended)
		})
	}
}

func TestBaseLaneConfig_ParseExtended(t *testing.T) {
	type args struct {
		tomlString string
		itemType   LaneTypes
		name       string
	}
	type testStruct struct {
		Host string  `pipelane:"host"`
		Port *string `pipelane:"port"`
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test extended host and port",
			args: args{
				tomlString: `
[input.input1]
buffer = 1
source_name = "int"
host = "0.0.0.0"
port = "8080"
`,
				itemType: InputType,
				name:     "input1",
			},
			want: testStruct{
				Host: "0.0.0.0",
				Port: stringPtr("8080"),
			},
		},
		{
			name: "test extended host port nil",
			args: args{
				tomlString: `
[input.input1]
buffer = 1
source_name = "int"
host = "0.0.0.0"
`,
				itemType: InputType,
				name:     "input1",
			},
			want: testStruct{
				Host: "0.0.0.0",
				Port: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tom, err := decodeTomlString(tt.args.tomlString)
			if err != nil {
				assert.Error(t, err)
				return
			}
			data, ok := tom[string(tt.args.itemType)].(map[string]any)
			if !ok {
				assert.Error(t, errors.New("not a map[string]any"))
				return
			}
			extended, ok := data[tt.args.name].(map[string]any)
			if !ok {
				assert.Error(t, errors.New("not a map[string]any"))
				return
			}
			cfg, err := NewBaseConfigWithTypeAndExtended(tt.args.itemType, tt.args.name, extended)
			if err != nil {
				assert.NotNil(t, err)
				return
			}
			var got testStruct
			err = cfg.ParseExtended(&got)
			if err != nil {
				assert.NotNil(t, err)
				return
			}
			assert.Equalf(t, got, tt.want, fmt.Sprintf("ParseExtended(%v)", tt.args.tomlString))
		})
	}
}

func TestBaseLaneConfig_ParseExtendedArrays(t *testing.T) {
	type args struct {
		tomlString string
		itemType   LaneTypes
		name       string
	}
	type testStruct struct {
		Hosts []string `pipelane:"hosts"`
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test array extended host and port",
			args: args{
				tomlString: `
[input.input1]
buffer = 1
source_name = "int"
hosts = ["0.0.0.0", "1.1.1.1"]
`,
				itemType: InputType,
				name:     "input1",
			},
			want: testStruct{
				Hosts: []string{"0.0.0.0", "1.1.1.1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tom, err := decodeTomlString(tt.args.tomlString)
			if err != nil {
				assert.Error(t, err)
				return
			}
			data, ok := tom[string(tt.args.itemType)].(map[string]any)
			if !ok {
				assert.Error(t, errors.New("not a map[string]any"))
				return
			}
			extended, ok := data[tt.args.name].(map[string]any)
			if !ok {
				assert.Error(t, errors.New("not a map[string]any"))
				return
			}
			cfg, err := NewBaseConfigWithTypeAndExtended(tt.args.itemType, tt.args.name, extended)
			if err != nil {
				assert.NotNil(t, err)
				return
			}
			var got testStruct
			err = cfg.ParseExtended(&got)
			if err != nil {
				assert.NotNil(t, err)
				return
			}
			assert.Equalf(t, got, tt.want, fmt.Sprintf("ParseExtended(%v)", tt.args.tomlString))
		})
	}
}
