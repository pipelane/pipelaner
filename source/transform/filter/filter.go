/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package remap

import (
	"encoding/json"
	"errors"

	"github.com/expr-lang/expr/vm"
	"github.com/rs/zerolog"

	"github.com/expr-lang/expr"

	"github.com/pipelane/pipelaner"
)

var (
	ErrInvalidDataType = errors.New("error invalid data type")
)

type EnvMap struct {
	Data map[string]any
}

type Config struct {
	Code string `pipelane:"code"`
}

type Filter struct {
	cfg     *pipelaner.BaseLaneConfig
	logger  zerolog.Logger
	program *vm.Program
}

func init() {
	pipelaner.RegisterMap("filter", &Filter{})
}

func (e *Filter) Init(ctx *pipelaner.Context) error {
	e.cfg = ctx.LaneItem().Config()
	e.logger = pipelaner.NewLogger()
	v := &Config{}
	err := e.cfg.ParseExtended(v)
	if err != nil {
		return err
	}

	program, err := expr.Compile(v.Code, expr.Env(EnvMap{}))
	if err != nil {
		return err
	}
	e.program = program
	return nil
}

func (e *Filter) Map(_ *pipelaner.Context, val any) any {
	var v map[string]any
	switch value := val.(type) {
	case map[string]any:
		v = value
	case string:
		b := []byte(value)
		err := json.Unmarshal(b, &v)
		if err != nil {
			return err
		}
	case []byte:
		err := json.Unmarshal(value, &v)
		if err != nil {
			return err
		}
	default:
		return ErrInvalidDataType
	}
	output, err := expr.Run(e.program, EnvMap{Data: v})
	if err != nil {
		e.logger.Err(err).Msg("Expr: output error")
		return err
	}
	if !output.(bool) {
		return nil
	}
	return val
}
