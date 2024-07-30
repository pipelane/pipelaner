package clickhouse

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/pipelane/pipelaner"
	"github.com/rs/zerolog"
)

type Clickhouse struct {
	logger      zerolog.Logger
	clickConfig *ClickhouseConfig
	client      *ClientClickhouse
}

func (c *Clickhouse) Init(ctx *pipelaner.Context) error {
	c.logger = pipelaner.NewLogger()
	err := ctx.LaneItem().Config().ParseExtended(c.clickConfig)
	if err != nil {
		return err
	}
	cli, err := NewClickhouseClient(ctx.Context(), c.clickConfig)
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
	sb.InsertInto(c.clickConfig.TableName).Cols(cols...).Values(values).SetFlavor(sqlbuilder.ClickHouse)

	sql, args := sb.Build()

	if _, err := c.client.Conn().Query(ctx, sql, args); err != nil {
		c.logger.Error().Err(err).Msgf("insert values clickhouse")
		return
	}
}

func (c *Clickhouse) Sink(ctx *pipelaner.Context, val any) {
	data := make(map[string]any)

	switch v := val.(type) {
	case json.RawMessage:
		if err := json.Unmarshal(v, &data); err != nil {
			c.logger.Error().Err(err).Msgf("unmarshal val")
			return
		}
	case []byte:
		if err := json.Unmarshal(v, &data); err != nil {
			c.logger.Error().Err(err).Msgf("unmarshal val")
			return
		}
	case map[string]any:
		data = v
	case chan any:
		for ch := range val.(chan any) {
			switch vv := ch.(type) {
			case json.RawMessage:
				if err := json.Unmarshal(vv, &data); err != nil {
					c.logger.Error().Err(err).Msgf("unmarshal value channel")
					return
				}
			case []byte:
				if err := json.Unmarshal(vv, &data); err != nil {
					c.logger.Error().Err(err).Msgf("unmarshal value channel")
					return
				}
			case map[string]any:
				data = vv
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
