/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package transform

import (
	"context"

	"github.com/antonmedv/expr/vm"
	"github.com/rs/zerolog"

	"github.com/antonmedv/expr"

	pipelane "github.com/pipelane/pipelaner"
)

type EnvMap struct {
	Data map[string]any
}

type ExprConfig struct {
	Code string `pipelane:"code"`
}

type Filter struct {
	cfg     *pipelane.BaseLaneConfig
	logger  zerolog.Logger
	program *vm.Program
}

func (e *Filter) Init(cfg *pipelane.BaseLaneConfig) error {
	e.cfg = cfg
	e.logger = pipelane.NewLogger()
	v := &ExprConfig{}
	err := cfg.ParseExtended(v)
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

func (e *Filter) Map(ctx context.Context, val any) any {
	output, err := expr.Run(e.program, EnvMap{Data: val.(map[string]any)})
	if err != nil {
		e.logger.Err(err).Msg("Expr: output error")
		return err
	}
	if !output.(bool) {
		return nil
	}
	return val
}
