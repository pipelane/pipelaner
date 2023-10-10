package pipelane

import "github.com/mitchellh/mapstructure"

type LaneTypes string

const (
	InputType LaneTypes = "input"
	LaneType  LaneTypes = "lane"
	SinkType  LaneTypes = "sink"
)

type config struct {
	Input map[string]any `pipeline:"input"`
	Lane  map[string]any `pipeline:"lane"`
	Sink  map[string]any `pipeline:"sink"`
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

type BaseConfig struct {
	BufferSize int64          `pipelane:"buffer"`
	SourceName string         `pipelane:"source_name"`
	Input      *string        `pipelane:"input"`
	Name       string         `pipelane:"-"`
	LaneType   LaneTypes      `pipelane:"-"`
	Extended   any            `pipelane:"-"`
	_extended  map[string]any `pipelane:"-"`
}

func NewBaseConfigWithTypeAndExtended(
	ItemType LaneTypes,
	extended map[string]any,
) (*BaseConfig, error) {
	c := BaseConfig{LaneType: ItemType, _extended: extended}
	err := decode(extended, &c)
	if err != nil {
		return nil, err
	}
	if c.BufferSize == 0 {
		c.BufferSize = 1
	}
	return &c, nil
}

func NewBaseConfig(val map[string]any) (*BaseConfig, error) {
	var cfg BaseConfig
	err := decode(val, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *BaseConfig) ParseExtended(v any) error {
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
