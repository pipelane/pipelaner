/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/mapstructure"
)

var (
	ErrNameNotBeEmptyString = errors.New("ErrNameNotBeEmptyString")
	ErrEnvNotFound          = errors.New("ErrEnvNotFound")
)

type LaneTypes string

const (
	InputType LaneTypes = "input"
	MapType   LaneTypes = "map"
	SinkType  LaneTypes = "sink"
)

type LogFormat string

const (
	LogFormatPlain LogFormat = "plain"
	LogFormatJSON  LogFormat = "json"
)

type logConfig struct {
	LogLevel       string    `pipelane:"log_level"`
	EnableConsole  bool      `pipelane:"log_enable_console"`
	EnableFile     bool      `pipelane:"log_enable_file"`
	FileDirectory  string    `pipelane:"log_file_directory"`
	FileName       string    `pipelane:"log_file_name"`
	FileMaxSize    int       `pipelane:"log_file_max_size"`
	FileMaxBackups int       `pipelane:"log_file_max_backups"`
	FileMaxAge     int       `pipelane:"log_file_max_age"`
	FileCompress   bool      `pipelane:"log_file_compress"`
	FileLocalTime  bool      `pipelane:"log_file_local_time"`
	LogFormat      LogFormat `pipelane:"log_format"`
}

type healthCheckConfig struct {
	HealthCheckHost   string `pipelane:"health_check_host"`
	HealthCheckPort   int    `pipelane:"health_check_port"`
	HealthCheckEnable bool   `pipelane:"health_check_enable"`
}

type metricsConfig struct {
	MetricsHost        string `pipelane:"metrics_host"`
	MetricsPort        int    `pipelane:"metrics_port"`
	MetricsServiceName string `pipelane:"metrics_service_name"`
	MetricsEnable      bool   `pipelane:"metrics_enable"`
}

type Config struct {
	logConfig         `pipelane:",squash"`
	healthCheckConfig `pipelane:",squash"`
	metricsConfig     `pipelane:",squash"`
	Input             map[string]any `pipeline:"input"`
	Map               map[string]any `pipeline:"map"`
	Sink              map[string]any `pipeline:"sink"`
}

type Internal struct {
	Name      string         `pipelane:"-"`
	LaneType  LaneTypes      `pipelane:"-"`
	Extended  any            `pipelane:"-"`
	_extended map[string]any `pipelane:"-"`
}

type BaseLaneConfig struct {
	OutputBufferSize           int64    `pipelane:"output_buffer"`
	Threads                    int64    `pipelane:"threads"`
	StartGCAfterMessageProcess bool     `pipelane:"start_gc_after_message_process"`
	SourceName                 string   `pipelane:"source_name"`
	Inputs                     []string `pipelane:"inputs"`
	Internal
}

func NewConfig(c map[string]any) (*Config, error) {
	cfg := &Config{}
	dC := &mapstructure.DecoderConfig{
		TagName: "pipelane",
		Result:  &cfg,
	}
	dec, err := mapstructure.NewDecoder(dC)
	if err != nil {
		return nil, err
	}
	err = dec.Decode(c)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func NewConfigFromFile(file string) (*Config, error) {
	c, err := ReadToml(file)
	if err != nil {
		return nil, err
	}
	cfg, err := NewConfig(c)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func NewBaseConfigWithTypeAndExtended(
	itemType LaneTypes,
	name string,
	extended map[string]any,
) (*BaseLaneConfig, error) {
	if name == "" {
		return nil, ErrNameNotBeEmptyString
	}
	c := BaseLaneConfig{
		Internal: Internal{
			LaneType:  itemType,
			Name:      name,
			_extended: extended,
		},
	}
	err := decode(extended, &c)
	if err != nil {
		return nil, err
	}
	if itemType == InputType {
		c.Inputs = nil
	}
	if c.OutputBufferSize == 0 {
		c.OutputBufferSize = int64(runtime.NumCPU())
	}
	if c.Threads == 0 {
		c.Threads = int64(runtime.NumCPU())
	}
	return &c, nil
}

func ReadToml(file string) (map[string]any, error) {
	var c map[string]any
	_, err := toml.DecodeFile(file, &c)
	if err != nil {
		return nil, err
	}
	err = recursiveReplace(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func recursiveReplace(cfg map[string]any) error {
	for k, v := range cfg {
		switch val := v.(type) {
		case map[string]any:
			return recursiveReplace(val)
		case string:
			e, err := findEnvValue(val)
			if err != nil {
				return err
			}
			cfg[k] = e
		case []any:
			if len(val) == 0 || len(val) > 1 {
				return fmt.Errorf("invalid env var %s array", k)
			}
			e, err := findEnvValue(val[0].(string))
			if err != nil {
				return err
			}
			cfg[k] = strings.Split(e, ",")
		}
	}
	return nil
}

func findEnvValue(val string) (string, error) {
	if strings.HasSuffix(val, "$") && strings.HasPrefix(val, "$") {
		envName := strings.ReplaceAll(val, "$", "")
		envName = strings.ReplaceAll(envName, "$", "")
		envName = strings.ReplaceAll(envName, " ", "")
		envValue := os.Getenv(strings.ToUpper(envName))
		if envValue == "" {
			return "", fmt.Errorf("env var %s not set", envValue)
		}
		return envValue, nil
	}
	return val, nil
}

func decodeTomlString(str string) (map[string]any, error) {
	var c map[string]any
	_, err := toml.Decode(str, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func NewBaseConfig(val map[string]any) (*BaseLaneConfig, error) {
	var cfg BaseLaneConfig
	err := decode(val, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *BaseLaneConfig) ParseExtended(v any) error {
	err := decode(c._extended, v)
	if err != nil {
		return err
	}
	c.Extended = v
	return nil
}

func decode(input map[string]any, output any) error {
	dC := &mapstructure.DecoderConfig{
		TagName: "pipelane",
		Result:  output,
		Squash:  true,
	}
	dec, err := mapstructure.NewDecoder(dC)
	if err != nil {
		return err
	}
	err = dec.Decode(input)
	if err != nil {
		return err
	}
	return nil
}
