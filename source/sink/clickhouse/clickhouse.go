package clickhouse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/pipelane/pipelaner"
)

type Clickhouse struct {
	logger      *zerolog.Logger
	clickConfig Config
	client      *LowLevelClickhouseClient
}

func init() {
	pipelaner.RegisterSink("clickhouse", &Clickhouse{})
}

func (c *Clickhouse) Init(ctx *pipelaner.Context) error {
	c.logger = ctx.Logger()

	err := ctx.LaneItem().Config().ParseExtended(&c.clickConfig)
	if err != nil {
		return err
	}
	cli, err := NewLowLevelClickhouseClient(ctx.Context(), c.clickConfig)
	if err != nil {
		return err
	}

	c.client = cli

	return nil
}

func AppendInput(
	input proto.Input,
	columnName string,
	data proto.Column,
) proto.Input {
	input = append(input, proto.InputColumn{Name: columnName, Data: data})
	return input
}

type column struct {
	str          *proto.ColStr
	flt          *proto.ColFloat64
	integer      *proto.ColInt64
	boolean      *proto.ColBool
	strArr       *proto.ColArr[string]
	intArr       *proto.ColArr[int64]
	fltArr       *proto.ColArr[float64]
	boolArr      *proto.ColArr[bool]
	arrStrArray  *proto.ColArr[[]string]
	arrIntArray  *proto.ColArr[[]int64]
	arrFltArray  *proto.ColArr[[]float64]
	arrBoolArray *proto.ColArr[[]bool]
	uid          *proto.ColUUID
	timestamp    *proto.ColDateTime64
}

// depending on the type of the column, write data to the proto column .

func (c *column) Append(v any) error {
	switch val := v.(type) {
	case string:
		if c.str != nil {
			c.str.Append(val)
		}
	case float64:
		if c.flt != nil {
			c.flt.Append(val)
		}
	case int64:
		if c.integer != nil {
			c.integer.Append(val)
		}
	case bool:
		if c.boolean != nil {
			c.boolean.Append(val)
		}
	case []string:
		if c.strArr != nil {
			c.strArr.Append(val)
		}
	case []float64:
		if c.fltArr != nil {
			c.fltArr.Append(val)
		}
	case []int64:
		if c.intArr != nil {
			c.intArr.Append(val)
		}
	case []bool:
		if c.boolArr != nil {
			c.boolArr.Append(val)
		}
	case [][]string:
		if c.arrStrArray != nil {
			c.arrStrArray.Append(val)
		}
	case [][]int64:
		if c.arrIntArray != nil {
			c.arrIntArray.Append(val)
		}
	case [][]float64:
		if c.arrFltArray != nil {
			c.arrFltArray.Append(val)
		}
	case [][]bool:
		if c.arrBoolArray != nil {
			c.arrBoolArray.Append(val)
		}
	case uuid.UUID:
		if c.uid != nil {
			c.uid.Append(val)
		}
	case time.Time:
		if c.timestamp != nil {
			c.timestamp.Append(val)
		}
	default:
		return errors.New("unknown type")
	}

	return nil
}

// buildProtoInput returns column map and input where column field
// is a pointer to input.Data, map key column name(input.Name)

func (c *Clickhouse) buildProtoInput(m map[string]any) (map[string]*column, proto.Input, error) {
	input := proto.Input{}

	columns := make(map[string]*column, len(m))

	for k, v := range m {
		col := new(column)
		switch v.(type) {
		case string:
			col.str = new(proto.ColStr)
			input = AppendInput(input, k, col.str)
		case float64:
			col.flt = new(proto.ColFloat64)
			input = AppendInput(input, k, col.flt)
		case int64:
			col.integer = new(proto.ColInt64)
			input = AppendInput(input, k, col.integer)
		case bool:
			col.boolean = new(proto.ColBool)
			input = AppendInput(input, k, col.boolean)
		case []string:
			col.strArr = new(proto.ColStr).Array()
			input = AppendInput(input, k, col.strArr)
		case []float64:
			col.fltArr = new(proto.ColFloat64).Array()
			input = AppendInput(input, k, col.fltArr)
		case []int64:
			col.intArr = new(proto.ColInt64).Array()
			input = AppendInput(input, k, col.intArr)
		case []bool:
			col.boolArr = new(proto.ColBool).Array()
			input = AppendInput(input, k, col.boolArr)
		case [][]string:
			col.arrStrArray = new(proto.ColArr[[]string])
			input = AppendInput(input, k, col.arrStrArray)
		case [][]int64:
			col.arrIntArray = new(proto.ColArr[[]int64])
			input = AppendInput(input, k, col.arrIntArray)
		case [][]float64:
			col.arrFltArray = new(proto.ColArr[[]float64])
			input = AppendInput(input, k, col.arrFltArray)
		case [][]bool:
			col.arrBoolArray = new(proto.ColArr[[]bool])
			input = AppendInput(input, k, col.arrBoolArray)
		case uuid.UUID:
			col.uid = new(proto.ColUUID)
			input = AppendInput(input, k, col.uid)
		case time.Time:
			col.timestamp = new(proto.ColDateTime64)
			input = AppendInput(input, k, col.timestamp)
		default:
			return nil, nil, fmt.Errorf("type val for column %s not found", k)
		}

		columns[k] = col
	}

	return columns, input, nil
}

func (c *Clickhouse) write(ctx context.Context, data []map[string]any) error {
	if len(data) == 0 {
		return fmt.Errorf("empty data")
	}

	columns, input, err := c.buildProtoInput(data[0])
	if err != nil {
		return fmt.Errorf("build proto input: %w", err)
	}

	conn, err := c.client.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("failed get conn: %w", err)
	}
	defer conn.Release()

	var blocks int
	if err = conn.Do(ctx, ch.Query{
		Body: input.Into(c.clickConfig.TableName),
		Settings: []ch.Setting{
			{
				Key:       "async_insert",
				Value:     c.clickConfig.AsyncInsert,
				Important: true,
			},
			{
				Key:       "wait_for_async_insert",
				Value:     c.clickConfig.WaitForAsyncInsert,
				Important: true,
			},
		},
		OnInput: func(_ context.Context) error {
			input.Reset()

			if blocks >= len(data)-1 {
				return io.EOF
			}
			for i := range data {
				for k, v := range data[i] {
					col, ok := columns[k]
					if !ok {
						return fmt.Errorf("column %s not found", k)
					}

					if err = col.Append(v); err != nil {
						return err
					}
				}

				blocks++
			}

			return nil
		},
		Input: input,
	}); err != nil {
		return fmt.Errorf("write batch: %w", err)
	}

	return nil
}

func (c *Clickhouse) Sink(ctx *pipelaner.Context, val any) {
	data := make(map[string]any)

	switch v := val.(type) {
	case json.RawMessage:
		if err := json.Unmarshal(v, &data); err != nil {
			c.logger.Error().Err(err).Msgf("RawMessage unmarshal")
			return
		}
	case []byte:
		if err := json.Unmarshal(v, &data); err != nil {
			c.logger.Error().Err(err).Msgf("[]byte unmarshal val")
			return
		}
	case map[string]any:
		data = v
	case chan any:
		listData := make([]map[string]any, 0, cap(v))

		for ch := range val.(chan any) {
			switch vv := ch.(type) {
			case json.RawMessage:
				if err := json.Unmarshal(vv, &data); err != nil {
					c.logger.Error().Err(err).Msgf("channel RawMessage unmarshal")
					return
				}
			case []byte:
				if err := json.Unmarshal(vv, &data); err != nil {
					c.logger.Error().Err(err).Msgf("channel []byte unmarshal")
					return
				}
			case map[string]any:
				data = vv
			default:
				c.logger.Error().Err(errors.New("unknown channel type"))
				return
			}

			listData = append(listData, data)
		}

		if err := c.write(ctx.Context(), listData); err != nil {
			c.logger.Error().Err(err).Msg("write")
			return
		}

		return
	default:
		c.logger.Error().Err(errors.New("unknown type val"))
		return
	}

	if err := c.write(ctx.Context(), []map[string]any{data}); err != nil {
		c.logger.Error().Err(err).Msg("write")
		return
	}
}
