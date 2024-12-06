/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/BurntSushi/toml"
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
		want      *Config
		wantError bool
	}{
		{
			name: "test inputs",
			args: args{
				tomlString: `
[input.input1]
output_buffer = 1
source_name = "int"
	`,
			},
			want: &Config{
				Input: map[string]any{
					"input1": map[string]any{
						"output_buffer": int64(1),
						"source_name":   "int",
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
output_buffer = 1
source_name = "int"
	`,
			},
			want: &Config{
				Input: nil,
				Map: map[string]any{
					"map2": map[string]any{
						"output_buffer": int64(1),
						"source_name":   "int",
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
output_buffer = 1
source_name = "int"
	`,
			},
			want: &Config{
				Input: nil,
				Map:   nil,
				Sink: map[string]any{
					"sink3": map[string]any{
						"output_buffer": int64(1),
						"source_name":   "int",
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
output_buffer = 10
source_name = "input_int"
[map.map2]
output_buffer = 20
source_name = "map_int"
[sink.sink3]
output_buffer = 30
source_name = "sink_int"
	`,
			},
			want: &Config{
				Input: map[string]any{
					"input1": map[string]any{
						"output_buffer": int64(10),
						"source_name":   "input_int",
					},
				},
				Map: map[string]any{
					"map2": map[string]any{
						"output_buffer": int64(20),
						"source_name":   "map_int",
					},
				},
				Sink: map[string]any{
					"sink3": map[string]any{
						"output_buffer": int64(30),
						"source_name":   "sink_int",
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
output_buffer = 10
source_name = "input_int"
[map.map2]
output_buffer = 20
source_name = "map_int"
[sink.sink3]
output_buffer = 30
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
			got, err := NewConfig(tom)
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
					"output_buffer": int64(1),
					"threads":       int64(1),
					"source_name":   "int",
					"host":          "0.0.0.0",
					"port":          "8080",
				},
			},
			want: &BaseLaneConfig{
				OutputBufferSize: 1,
				Threads:          1,
				SourceName:       "int",
				Inputs:           nil,
				Internal: Internal{
					Name:     "input1",
					LaneType: InputType,
					Extended: nil,
					_extended: map[string]any{
						"output_buffer": int64(1),
						"threads":       int64(1),
						"source_name":   "int",
						"host":          "0.0.0.0",
						"port":          "8080",
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
					"output_buffer": int64(1),
					"threads":       int64(1),
					"source_name":   "int",
					"host":          "0.0.0.0",
					"port":          "8080",
				},
			},
			want: &BaseLaneConfig{
				OutputBufferSize: 1,
				Threads:          1,
				SourceName:       "int",
				Inputs:           nil,
				Internal: Internal{
					Name:     "input1",
					LaneType: InputType,
					Extended: nil,
					_extended: map[string]any{
						"output_buffer": int64(1),
						"threads":       int64(1),
						"source_name":   "int",
						"host":          "0.0.0.0",
						"port":          "8080",
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
					"output_buffer": int64(runtime.NumCPU()),
					"threads":       int64(1),
					"source_name":   "int",
					"host":          "0.0.0.0",
					"port":          "8080",
				},
			},
			want: &BaseLaneConfig{
				OutputBufferSize: int64(runtime.NumCPU()),
				Threads:          1,
				SourceName:       "int",
				Inputs:           nil,
				Internal: Internal{
					Name:     "input1",
					LaneType: InputType,
					Extended: nil,
					_extended: map[string]any{
						"output_buffer": int64(runtime.NumCPU()),
						"threads":       int64(1),
						"source_name":   "int",
						"host":          "0.0.0.0",
						"port":          "8080",
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
					"inputs":        []string{"input_name"},
					"source_name":   "int",
					"host":          "0.0.0.0",
					"port":          "8080",
					"output_buffer": int64(1),
					"threads":       int64(1),
				},
			},
			want: &BaseLaneConfig{
				OutputBufferSize: 1,
				Threads:          1,
				SourceName:       "int",
				Inputs:           []string{"input_name"},
				Internal: Internal{
					Name:     "sink1",
					LaneType: SinkType,
					Extended: nil,
					_extended: map[string]any{
						"source_name":   "int",
						"host":          "0.0.0.0",
						"port":          "8080",
						"inputs":        []string{"input_name"},
						"output_buffer": int64(1),
						"threads":       int64(1),
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
					"inputs":        []string{"input_name"},
					"source_name":   "int",
					"host":          "0.0.0.0",
					"port":          "8080",
					"output_buffer": int64(1),
					"threads":       int64(1),
				},
			},
			want: &BaseLaneConfig{
				OutputBufferSize: 1,
				Threads:          1,
				SourceName:       "int",
				Inputs:           []string{"input_name"},
				Internal: Internal{
					Name:     "map1",
					LaneType: MapType,
					Extended: nil,
					_extended: map[string]any{
						"output_buffer": int64(1),
						"threads":       int64(1),
						"source_name":   "int",
						"host":          "0.0.0.0",
						"port":          "8080",
						"inputs":        []string{"input_name"},
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
output_buffer = 1
threads = 1
source_name = "int"
host = "0.0.0.0"
port = "8080"
`,
				itemType: InputType,
				name:     "input1",
			},
			want: &BaseLaneConfig{
				OutputBufferSize: 1,
				Threads:          1,
				SourceName:       "int",
				Inputs:           nil,
				Internal: Internal{
					Name:     "input1",
					LaneType: InputType,
					Extended: nil,
					_extended: map[string]any{
						"output_buffer": int64(1),
						"threads":       int64(1),
						"source_name":   "int",
						"host":          "0.0.0.0",
						"port":          "8080",
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
threads = 1
output_buffer = 1
`,
				itemType: InputType,
				name:     "input1",
			},
			want: &BaseLaneConfig{
				OutputBufferSize: 1,
				Threads:          1,
				SourceName:       "int",
				Inputs:           nil,
				Internal: Internal{
					Name:     "input1",
					LaneType: InputType,
					Extended: nil,
					_extended: map[string]any{
						"source_name":   "int",
						"host":          "0.0.0.0",
						"port":          "8080",
						"output_buffer": int64(1),
						"threads":       int64(1),
					},
				},
			},
		},
		{
			name: "test input type with buffer zero",
			args: args{
				tomlString: `
[input.input1]
output_buffer = 1
threads = 1
source_name = "int"
host = "0.0.0.0"
port = "8080"
`,
				itemType: InputType,
				name:     "input1",
			},
			want: &BaseLaneConfig{
				OutputBufferSize: 1,
				Threads:          1,
				SourceName:       "int",
				Inputs:           nil,
				Internal: Internal{
					Name:     "input1",
					LaneType: InputType,
					Extended: nil,
					_extended: map[string]any{
						"source_name":   "int",
						"host":          "0.0.0.0",
						"port":          "8080",
						"output_buffer": int64(1),
						"threads":       int64(1),
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
output_buffer = 1
threads = 1
`,
				itemType: SinkType,
				name:     "sink1",
			},
			want: &BaseLaneConfig{
				OutputBufferSize: 1,
				Threads:          1,
				SourceName:       "int",
				Inputs:           []string{"input_name"},
				Internal: Internal{
					Name:     "sink1",
					LaneType: SinkType,
					Extended: nil,
					_extended: map[string]any{
						"source_name":   "int",
						"host":          "0.0.0.0",
						"port":          "8080",
						"inputs":        []any{"input_name"},
						"output_buffer": int64(1),
						"threads":       int64(1),
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
output_buffer = 1
threads = 1
`,
				itemType: MapType,
				name:     "map1",
			},
			want: &BaseLaneConfig{
				OutputBufferSize: 1,
				Threads:          1,
				SourceName:       "int",
				Inputs:           []string{"input_name"},
				Internal: Internal{
					Name:     "map1",
					LaneType: MapType,
					Extended: nil,
					_extended: map[string]any{
						"source_name":   "int",
						"host":          "0.0.0.0",
						"port":          "8080",
						"inputs":        []any{"input_name"},
						"output_buffer": int64(1),
						"threads":       int64(1),
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
output_buffer = 1
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
output_buffer = 1
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
output_buffer = 1
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

func Test_injectEnvs(t *testing.T) {
	type args struct {
		cfg string
	}
	tests := []struct {
		name      string
		args      args
		want      map[string]any
		wantError bool
		setup     func()
	}{
		{
			name: "test inject single string value",
			args: args{
				cfg: `
# Get normalized data
[input.kafka_consumer]
source_name = "kafka"
batch_size = "32KiB"
sasl_enabled = true
sasl_mechanism = "$KAFKA_SASL_MECHANISM$"
sasl_password = "$KAFKA_SASL_PASSWORD$"
sasl_username = "$KAFKA_SASL_USERNAME$"
`,
			},
			want: map[string]any{
				"input": map[string]any{
					"kafka_consumer": map[string]any{
						"source_name":    "kafka",
						"batch_size":     "32KiB",
						"sasl_enabled":   true,
						"sasl_mechanism": "PLAIN",
						"sasl_password":  "123",
						"sasl_username":  "321",
					},
				},
			},
			wantError: false,
			setup: func() {
				err := os.Setenv("KAFKA_SASL_MECHANISM", "PLAIN")
				assert.NoError(t, err)
				err = os.Setenv("KAFKA_SASL_PASSWORD", "123")
				assert.NoError(t, err)
				err = os.Setenv("KAFKA_SASL_USERNAME", "321")
				assert.NoError(t, err)
			},
		},
		{
			name: "test inject single lowercased string value",
			args: args{
				cfg: `
# Get normalized data
[input.kafka_consumer]
source_name = "kafka"
batch_size = "32KiB"
sasl_enabled = true
sasl_mechanism = "$kafka_sasl_mechanism$"
sasl_password = "$kafka_sasl_password$"
sasl_username = "$kafka_sasl_username$"
`,
			},
			want: map[string]any{
				"input": map[string]any{
					"kafka_consumer": map[string]any{
						"source_name":    "kafka",
						"batch_size":     "32KiB",
						"sasl_enabled":   true,
						"sasl_mechanism": "PLAIN",
						"sasl_password":  "123",
						"sasl_username":  "321",
					},
				},
			},
			wantError: false,
			setup: func() {
				err := os.Setenv("KAFKA_SASL_MECHANISM", "PLAIN")
				assert.NoError(t, err)
				err = os.Setenv("KAFKA_SASL_PASSWORD", "123")
				assert.NoError(t, err)
				err = os.Setenv("KAFKA_SASL_USERNAME", "321")
				assert.NoError(t, err)
			},
		},
		{
			name: "test inject array and single string value",
			args: args{
				cfg: `
# Get normalized data
[input.kafka_consumer]
source_name = "kafka"
batch_size = "32KiB"
sasl_enabled = true
sasl_mechanism = "$KAFKA_SASL_MECHANISM$"
sasl_password = "$KAFKA_SASL_PASSWORD$"
sasl_username = "$KAFKA_SASL_USERNAME$"
topics = ["$KAFKA_TOPICS$"]
`,
			},
			want: map[string]any{
				"input": map[string]any{
					"kafka_consumer": map[string]any{
						"source_name":    "kafka",
						"batch_size":     "32KiB",
						"sasl_enabled":   true,
						"sasl_mechanism": "PLAIN",
						"sasl_password":  "123",
						"sasl_username":  "321",
						"topics":         []any{"1", "2"},
					},
				},
			},
			wantError: false,
			setup: func() {
				err := os.Setenv("KAFKA_SASL_MECHANISM", "PLAIN")
				assert.NoError(t, err)
				err = os.Setenv("KAFKA_SASL_PASSWORD", "123")
				assert.NoError(t, err)
				err = os.Setenv("KAFKA_SASL_USERNAME", "321")
				assert.NoError(t, err)
				err = os.Setenv("KAFKA_TOPICS", "1,2")
				assert.NoError(t, err)
			},
		},
		{
			name: "test inject single string with error",
			args: args{
				cfg: `
# Get normalized data
[input.kafka_consumer]
source_name = "kafka"
batch_size = "32KiB"
sasl_enabled = true
sasl_mechanism = "$KAFKA_SASL_MECHANISM$"
sasl_password = "$KAFKA_SASL_PASSWORD$"
sasl_username = "$KAFKA_SASL_USERNAME$"
`,
			},
			want: map[string]any{
				"input": map[string]any{
					"kafka_consumer": map[string]any{
						"source_name":    "kafka",
						"batch_size":     "32KiB",
						"sasl_enabled":   true,
						"sasl_mechanism": "PLAIN",
						"sasl_password":  "123",
						"sasl_username":  "321",
					},
				},
			},
			wantError: true,
			setup: func() {
				err := os.Setenv("KAFKA_SASL_MECHANISM", "PLAIN")
				assert.NoError(t, err)
				err = os.Setenv("KAFKA_SASL_PASSWORD", "123")
				assert.NoError(t, err)
			},
		},
		{
			name: "test inject array of int",
			args: args{
				cfg: `
# Get normalized data
[input.kafka_consumer]
source_name = "kafka"
batch_size = "32KiB"
sasl_enabled = true
sasl_mechanism = "$KAFKA_SASL_MECHANISM$"
sasl_password = "$KAFKA_SASL_PASSWORD$"
sasl_username = "$KAFKA_SASL_USERNAME$"
slice_array = [1, 2, 3] 
`,
			},
			want: map[string]any{
				"input": map[string]any{
					"kafka_consumer": map[string]any{
						"source_name":    "kafka",
						"batch_size":     "32KiB",
						"sasl_enabled":   true,
						"sasl_mechanism": "PLAIN",
						"sasl_password":  "123",
						"sasl_username":  "321",
						"slice_array":    []any{int64(1), int64(2), int64(3)},
					},
				},
			},
			wantError: false,
			setup: func() {
				err := os.Setenv("KAFKA_SASL_MECHANISM", "PLAIN")
				assert.NoError(t, err)
				err = os.Setenv("KAFKA_SASL_PASSWORD", "123")
				assert.NoError(t, err)
				err = os.Setenv("KAFKA_SASL_USERNAME", "321")
				assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			var got map[string]any
			_, err := toml.Decode(tt.args.cfg, &got)
			assert.NoError(t, err)
			got, err = recursiveReplace(got)
			if tt.wantError && err != nil {
				assert.Error(t, err, "injectEnvs() error = %v, wantErr %v", err)
				return
			}
			assert.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("injectEnvs()\n got = %v,\n want %v", got, tt.want)
			}
		})
	}
}
