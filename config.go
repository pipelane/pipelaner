/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"errors"
	"fmt"
	"reflect"
	"time"

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

type KafkaConfig struct {
	KafkaBrokers           string   `pipelane:"brokers"`
	KafkaVersion           string   `pipelane:"version"`
	KafkaOffsetNewest      bool     `pipelane:"offset_newest"`
	KafkaSASLEnabled       bool     `pipelane:"sasl_enabled"`
	KafkaSASLMechanism     string   `pipelane:"sasl_mechanism"`
	KafkaSASLUsername      string   `pipelane:"sasl_username"`
	KafkaSASLPassword      string   `pipelane:"sasl_password"`
	KafkaAutoCommitEnabled bool     `pipelane:"auto_commit_enabled"`
	KafkaConsumerGroupId   string   `pipelane:"consumer_group_id"`
	KafkaTopics            []string `pipelane:"topics"`
	KafkaAutoOffsetReset   string   `pipelane:"auto_offset_reset"`
	KafkaBatchSize         int      `pipelane:"batch_size"`
	KafkaSchemaRegistry    string   `pipelane:"schema_registry"`
	Internal
}

type ClickHouseConfig struct {
	Address                  string        `pipelane:"address"`
	User                     string        `pipelane:"user"`
	Password                 string        `pipelane:"password"`
	Database                 string        `pipelane:"database"`
	MigrationEngine          string        `pipelane:"migration_engine"`
	MigrationsPathClickhouse string        `pipelane:"migrations_path_clickhouse"`
	MaxExecutionTime         time.Duration `pipelane:"max_execution_time"`
	ConnMaxLifetime          time.Duration `pipelane:"conn_max_lifetime"`
	DialTimeout              time.Duration `pipelane:"dial_timeout"`
	MaxOpenConns             int           `pipelane:"max_open_conns"`
	MaxIdleConns             int           `pipelane:"max_idle_conns"`
	BlockBufferSize          uint8         `pipelane:"block_buffer_size"`
	MaxCompressionBuffer     string        `pipelane:"max_compression_buffer"`
	EnableDebug              bool          `pipelane:"enable_debug"`
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

func NewKafkaConfig(val map[string]any) (*KafkaConfig, error) {
	var cfg KafkaConfig
	err := decode(val, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func NewClickHouseConfig(val map[string]any) (*ClickHouseConfig, error) {
	var cfg ClickHouseConfig
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

func CastConfig[K, V any](config K) *V {
	if !(reflect.TypeOf(config).Kind() == reflect.Ptr && reflect.TypeOf(config).Elem().Kind() == reflect.Struct) {
		panic("Config is not struct")
	}

	res := new(V)
	if reflect.ValueOf(res).Kind() == reflect.Struct {
		panic("V is not struct")
	}

	val := reflect.ValueOf(config).Elem()
	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name

		valueField := val.FieldByName(fieldName)
		if !valueField.IsValid() {
			panic(fmt.Sprintf("No such field: %s in valueField", fieldName))
		}

		setField := reflect.ValueOf(res).Elem().FieldByName(fieldName)
		if !setField.IsValid() {
			panic(fmt.Sprintf("No such field: %s in setField", fieldName))
		}

		if valueField.Type() != setField.Type() {
			panic(fmt.Sprintf("Cannot cast %s to %s", valueField.Type(), setField.Type()))
		}

		if !setField.CanSet() {
			panic(fmt.Sprintf("Cannot set %s field value", fieldName))
		}

		setField.Set(valueField)
	}

	return res
}
