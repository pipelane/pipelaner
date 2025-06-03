/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package clickhouse

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"github.com/google/uuid"
	"github.com/pipelane/pipelaner/gen/source/sink"
	"github.com/pipelane/pipelaner/pipeline/components"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/shopspring/decimal"
)

func init() {
	source.RegisterSink("clickhouse", &Clickhouse{})
}

type Clickhouse struct {
	components.Logger
	clickConfig sink.Clickhouse
	client      *Client
}

func (c *Clickhouse) Init(cfg sink.Sink) error {
	clickCfg, ok := cfg.(sink.Clickhouse)
	if !ok {
		return fmt.Errorf("invalid clickhouse config type: %T", cfg)
	}
	cli, err := NewClickhouseClient(context.Background(), clickCfg)
	if err != nil {
		return fmt.Errorf("init clickhouse client: %w", err)
	}
	c.client = cli
	c.clickConfig = clickCfg
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
	dcml         *proto.ColDecimal256
	boolean      *proto.ColBool
	strArr       *proto.ColArr[string]
	intArr       *proto.ColArr[int64]
	dcmlArr      *proto.ColArr[proto.Decimal256]
	fltArr       *proto.ColArr[float64]
	boolArr      *proto.ColArr[bool]
	arrStrArray  *proto.ColArr[[]string]
	arrIntArray  *proto.ColArr[[]int64]
	arrDcmlArray *proto.ColArr[[]proto.Decimal256]
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
	case decimal.Decimal:
		if c.dcml != nil {
			decimal256, err := ConvertToDecimal256(val, 10) // TODO: col scale
			if err != nil {
				return err
			}
			c.dcml.Append(decimal256)
		}
	case []decimal.Decimal:
		if c.dcmlArr != nil {
			arr := make([]proto.Decimal256, 0, len(val))
			for _, d := range val {
				decimal256, err := ConvertToDecimal256(d, 10) // TODO: col scale
				if err != nil {
					return err
				}
				arr = append(arr, decimal256)
			}
			c.dcmlArr.Append(arr)
		}
	case [][]decimal.Decimal:
		if c.arrDcmlArray != nil {
			doubleArr := make([][]proto.Decimal256, 0, len(val))
			for _, arr := range val {
				nestedArr := make([]proto.Decimal256, 0, len(arr))
				for _, d := range arr {
					decimal256, err := ConvertToDecimal256(d, 10)
					if err != nil {
						return err
					}
					nestedArr = append(nestedArr, decimal256)
				}

				doubleArr = append(doubleArr, nestedArr)
			}

			c.arrDcmlArray.Append(doubleArr)
		}
	default:
		return errors.New("unknown type")
	}

	return nil
}

// buildProtoInput returns column transform and input where column field
// is a pointer to input.Data, transform key column name(input.Name)

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
		case decimal.Decimal:
			col.dcml = new(proto.ColDecimal256)
			input = append(input, proto.InputColumn{Name: k, Data: proto.Alias(col.dcml, "Decimal(76, 10)")})
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
		case []decimal.Decimal:
			col.dcmlArr = proto.NewArray[proto.Decimal256](new(proto.ColDecimal256))
			input = append(input, proto.InputColumn{Name: k, Data: proto.Alias(col.dcmlArr, "Array(Decimal(76, 10))")})
		case []bool:
			col.boolArr = proto.NewArray[bool](new(proto.ColBool))
			input = AppendArrayInput(input, k, col.boolArr)
		case [][]string:
			col.arrStrArray = proto.NewArray[[]string](proto.NewArray[string](new(proto.ColStr)))
			input = AppendArrayInput(input, k, col.arrStrArray)
		case [][]int64:
			col.arrIntArray = proto.NewArray[[]int64](proto.NewArray[int64](new(proto.ColInt64)))
			input = AppendArrayInput(input, k, col.arrIntArray)
		case [][]decimal.Decimal:
			col.arrDcmlArray = proto.NewArray[[]proto.Decimal256](proto.NewArray[proto.Decimal256](new(proto.ColDecimal256)))
			input = append(input, proto.InputColumn{Name: k, Data: proto.Alias(col.arrDcmlArray, "Array(Array(Decimal(76, 10)))")})
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

func ConvertToDecimal256(v decimal.Decimal, targetScale int32) (proto.Decimal256, error) {
	var bi *big.Int
	bi = decimal.NewFromBigInt(v.Coefficient(), v.Exponent()+targetScale).BigInt()
	dest := make([]byte, 32)
	bigIntToRaw(dest, bi)
	return proto.Decimal256{
		Low: proto.UInt128{
			Low:  binary.LittleEndian.Uint64(dest[0 : 64/8]),
			High: binary.LittleEndian.Uint64(dest[64/8 : 128/8]),
		},
		High: proto.UInt128{
			Low:  binary.LittleEndian.Uint64(dest[128/8 : 192/8]),
			High: binary.LittleEndian.Uint64(dest[192/8 : 256/8]),
		},
	}, nil
}

func bigIntToRaw(dest []byte, v *big.Int) {
	var sign int
	if v.Sign() < 0 {
		v.Not(v).FillBytes(dest)
		sign = -1
	} else {
		v.FillBytes(dest)
	}
	endianSwap(dest, sign < 0)
}

func endianSwap(src []byte, not bool) {
	for i := 0; i < len(src)/2; i++ {
		if not {
			src[i], src[len(src)-i-1] = ^src[len(src)-i-1], ^src[i]
		} else {
			src[i], src[len(src)-i-1] = src[len(src)-i-1], src[i]
		}
	}
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
	}
	shouldEnd := false
	query := ch.Query{
		Body: input.Into(c.clickConfig.GetTableName()),
		Settings: []ch.Setting{
			{
				Key:       "async_insert",
				Value:     c.clickConfig.GetAsyncInsert(),
				Important: true,
			},
			{
				Key:       "wait_for_async_insert",
				Value:     c.clickConfig.GetWaitForAsyncInsert(),
				Important: true,
			},
			{
				Key:   "max_partitions_per_insert_block",
				Value: fmt.Sprintf("%d", c.clickConfig.GetMaxPartitionsPerInsertBlock()),
			},
		},
		OnInput: func(_ context.Context) error {
			input.Reset()
			if shouldEnd {
				return io.EOF
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
			newData, ok := <-chData
			if !ok && newData == nil {
				shouldEnd = true
				return nil
			}
			data, err = c.getMap(newData)
			if err != nil {
				return err
			}

			return nil
		},
		Input: input,
	}
	if err = conn.Do(ctx, query); err != nil {
		return fmt.Errorf("write batch: %w", err)
	}
	return nil
}

func (c *Clickhouse) Sink(val any) {
	data := make(map[string]any)
	var chData chan any
	switch v := val.(type) {
	case json.RawMessage:
		chData = make(chan any, 1)
		if err := json.Unmarshal(v, &data); err != nil {
			c.Log().Error().Err(err).Msgf("channel RawMessage unmarshal")
			return
		}
		chData <- data
		close(chData)
	case []byte:
		chData = make(chan any, 1)
		if err := json.Unmarshal(v, &data); err != nil {
			c.Log().Error().Err(err).Msgf("channel []byte unmarshal")
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
		c.Log().Error().Err(errors.New("unknown type val")).Msg("failed write clickhouse")
		return
	}
	if err := c.write(context.Background(), chData); err != nil {
		c.Log().Error().Err(err).Msg("write")
	}
}
