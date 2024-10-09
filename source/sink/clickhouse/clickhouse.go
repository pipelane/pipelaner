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
	var clickConfig Config
	err := ctx.LaneItem().Config().ParseExtended(&clickConfig)
	if err != nil {
		return err
	}
	if clickConfig.AsyncInsert == "" {
		clickConfig.AsyncInsert = "1"
	}
	if clickConfig.WaitForAsyncInsert == "" {
		clickConfig.WaitForAsyncInsert = "1"
	}
	c.clickConfig = clickConfig
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

func AppendArrayInput[T any](
	input proto.Input,
	columnName string,
	data *proto.ColArr[T],
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
			col.strArr = proto.NewArray[string](new(proto.ColStr))
			input = AppendArrayInput(input, k, col.strArr)
		case []float64:
			col.fltArr = proto.NewArray[float64](new(proto.ColFloat64))
			input = AppendArrayInput(input, k, col.fltArr)
		case []int64:
			col.intArr = proto.NewArray[int64](new(proto.ColInt64))
			input = AppendArrayInput(input, k, col.intArr)
		case []bool:
			col.boolArr = proto.NewArray[bool](new(proto.ColBool))
			input = AppendArrayInput(input, k, col.boolArr)
		case [][]string:
			col.arrStrArray = proto.NewArray[[]string](proto.NewArray[string](new(proto.ColStr)))
			input = AppendArrayInput(input, k, col.arrStrArray)
		case [][]int64:
			col.arrIntArray = proto.NewArray[[]int64](proto.NewArray[int64](new(proto.ColInt64)))
			input = AppendArrayInput(input, k, col.arrIntArray)
		case [][]float64:
			col.arrFltArray = proto.NewArray[[]float64](proto.NewArray[float64](new(proto.ColFloat64)))
			input = AppendArrayInput(input, k, col.arrFltArray)
		case [][]bool:
			col.arrBoolArray = proto.NewArray[[]bool](proto.NewArray[bool](new(proto.ColBool)))
			input = AppendArrayInput(input, k, col.arrBoolArray)
		case uuid.UUID:
			col.uid = new(proto.ColUUID)
			input = AppendInput(input, k, col.uid)
		case time.Time:
			col.timestamp = new(proto.ColDateTime64).WithPrecision(proto.PrecisionMicro)
			input = AppendInput(input, k, col.timestamp)
		default:
			return nil, nil, fmt.Errorf("type val for column %s not found", k)
		}

		columns[k] = col
	}

	return columns, input, nil
}

func (c *Clickhouse) getMap(val any) (map[string]any, error) {
	var d map[string]any

	switch v := val.(type) {
	case json.RawMessage:
		if err := json.Unmarshal(v, &d); err != nil {
			return nil, fmt.Errorf("RawMessage unmarshal")
		}
	case []byte:
		if err := json.Unmarshal(v, &d); err != nil {
			return nil, fmt.Errorf("channel []byte unmarshal")
		}
	case map[string]any:
		d = v
	}

	return d, nil
}

func (c *Clickhouse) write(ctx context.Context, chData chan any) error {
	conn, err := c.client.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("failed get conn: %w", err)
	}
	defer conn.Release()

	var (
		input   proto.Input
		columns map[string]*column
		data    map[string]any
	)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case v, ok := <-chData:
		if !ok && v == nil {
			return nil
		}
		data, err = c.getMap(v)
		if err != nil {
			return err
		}

		columns, input, err = c.buildProtoInput(data)
		if err != nil {
			return fmt.Errorf("build proto input: %w", err)
		}

		for k, val := range data {
			col, okC := columns[k]
			if !okC {
				return fmt.Errorf("column %s not found", k)
			}
			if err = col.Append(val); err != nil {
				return err
			}
		}
	}

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
			isClose := true

			for newData := range chData {
				isClose = false
				data, err = c.getMap(newData)
				if err != nil {
					return err
				}
				for k, v := range data {
					col, okC := columns[k]
					if !okC {
						return fmt.Errorf("column %s not found", k)
					}
					if err = col.Append(v); err != nil {
						return err
					}
				}
				return nil
			}
			if isClose {
				return io.EOF
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
	var chData chan any
	switch v := val.(type) {
	case json.RawMessage:
		chData = make(chan any, 1)
		if err := json.Unmarshal(v, &data); err != nil {
			c.logger.Error().Err(err).Msgf("channel RawMessage unmarshal")
			return
		}
		chData <- data
		close(chData)
	case []byte:
		chData = make(chan any, 1)
		if err := json.Unmarshal(v, &data); err != nil {
			c.logger.Error().Err(err).Msgf("channel []byte unmarshal")
			return
		}
		chData <- data
		close(chData)
	case map[string]any:
		chData = make(chan any, 1)
		chData <- v
		close(chData)
	case chan any:
		chData = v
	default:
		c.logger.Error().Err(errors.New("unknown type val"))
		return
	}
	if err := c.write(ctx.Context(), chData); err != nil {
		c.logger.Error().Err(err).Msg("write")
	}
}
