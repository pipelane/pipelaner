/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package remap

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/pipeline/source"
)

func init() {
	source.RegisterTransform("remap", &Remap{})
}

var (
	ErrInvalidDataType = errors.New("error invalid data type")
)

type EnvMap struct {
	Data any
}

type Config struct {
	Code string `pipelane:"code"`
}

type Remap struct {
	program *vm.Program
}

func (e *Remap) Init(cfg transform.Transform) error {
	rCfg, ok := cfg.(transform.Remap)
	if !ok {
		return fmt.Errorf("invalid remap config type: %T", cfg)
	}
	program, err := expr.Compile(rCfg.GetCode(), expr.Env(EnvMap{}))
	if err != nil {
		return err
	}
	e.program = program
	return nil
}

func (e *Remap) Transform(val any) any {
	var v any
	switch value := val.(type) {
	case map[string]any:
		v = value
	case map[string][]any:
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
		return err
	}
	return output
}
