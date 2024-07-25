package clickhouse

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/pipelane/go-kit/clickhouse"
	"github.com/pipelane/go-kit/config"
	"github.com/pipelane/pipelaner"
	"github.com/rs/zerolog"
)

type Clickhouse struct {
	logger zerolog.Logger
	cfg    *pipelaner.ClickHouseConfig
	client *clickhouse.ClientClickhouse
	table  string
}

func NewClickhouse(logger zerolog.Logger, cfg *pipelaner.ClickHouseConfig, table string) *Clickhouse {
	return &Clickhouse{
		logger: logger,
		cfg:    cfg,
		table:  table,
	}
}

func (c *Clickhouse) Init(ctx *pipelaner.Context) error {
	c.logger = pipelaner.NewLogger()

	castCfg := pipelaner.CastConfig[*pipelaner.ClickHouseConfig, config.ClickHouse](c.cfg)

	cli, err := clickhouse.NewClickhouseClient(ctx.Context(), castCfg)
	if err != nil {
		return err
	}

	c.client = cli

	return nil
}

func (c *Clickhouse) write(ctx context.Context, data map[string]any) {
	cols := make([]string, 0, len(data))
	values := make([]any, 0, len(data))

	for k, v := range data {
		cols = append(cols, k)
		values = append(values, v)
	}

	sb := sqlbuilder.NewInsertBuilder()
	sb.InsertInto(c.table).Cols(cols...).Values(values).SetFlavor(sqlbuilder.ClickHouse)

	sql, args := sb.Build()

	if _, err := c.client.Conn().Query(ctx, sql, args); err != nil {
		c.logger.Error().Err(err).Msgf("insert values clickhouse")
		return
	}
}

func (c *Clickhouse) Sink(ctx *pipelaner.Context, val any) {
	data := make(map[string]any)

	switch val.(type) {
	case json.RawMessage, []byte:
		if err := json.Unmarshal(val.([]byte), &data); err != nil {
			c.logger.Error().Err(err).Msgf("unmarshal val")
			return
		}
	case map[string]any:
		data = val.(map[string]any)
	case chan any:
		for ch := range val.(chan any) {
			switch ch.(type) {
			case json.RawMessage, []byte:
				if err := json.Unmarshal(ch.([]byte), &data); err != nil {
					c.logger.Error().Err(err).Msgf("unmarshal value channel")
					return
				}
			case map[string]any:
				data = ch.(map[string]any)
			default:
				c.logger.Error().Err(errors.New("unknown type channel"))
			}

			c.write(ctx.Context(), data)
		}

		return
	default:
		c.logger.Error().Err(errors.New("unknown type val"))
		return
	}

	c.write(ctx.Context(), data)
}
