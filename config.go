/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/mapstructure"
)

var (
	ErrNameNotBeEmptyString = errors.New("ErrNameNotBeEmptyString")
)

type LaneTypes string

const (
	InputType LaneTypes = "input"
	MapType   LaneTypes = "map"
	SinkType  LaneTypes = "sink"
)

type logConfig struct {
	LogLevel       string `pipelane:"log_level"`
	EnableConsole  bool   `pipelane:"log_enable_console"`
	EnableFile     bool   `pipelane:"log_enable_file"`
	FileDirectory  string `pipelane:"log_file_directory"`
	FileName       string `pipelane:"log_file_name"`
	FileMaxSize    int    `pipelane:"log_file_max_size"`
	FileMaxBackups int    `pipelane:"log_file_max_backups"`
	FileMaxAge     int    `pipelane:"log_file_max_age"`
	FileCompress   bool   `pipelane:"log_file_compress"`
	FileLocalTime  bool   `pipelane:"log_file_local_time"`
}

type config struct {
	logConfig `pipelane:",squash"`
	Input     map[string]any `pipeline:"input"`
	Map       map[string]any `pipeline:"map"`
	Sink      map[string]any `pipeline:"sink"`
}

type Internal struct {
	Name      string         `pipelane:"-"`
	LaneType  LaneTypes      `pipelane:"-"`
	Extended  any            `pipelane:"-"`
	_extended map[string]any `pipelane:"-"`
}

type BaseLaneConfig struct {
	BufferSize int64    `pipelane:"buffer"`
	Threads    *int64   `pipelane:"threads"`
	SourceName string   `pipelane:"source_name"`
	Inputs     []string `pipelane:"inputs"`
	Internal
}

func newConfig(c map[string]any) (*config, error) {
	cfg := &config{}
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
	if c.BufferSize == 0 {
		c.BufferSize = 1
	}
	return &c, nil
}

func readToml(file string) (map[string]any, error) {
	var c map[string]any
	_, err := toml.DecodeFile(file, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
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
