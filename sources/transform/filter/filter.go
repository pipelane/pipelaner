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
	"github.com/pipelane/pipelaner/pipeline/node"
	"github.com/pipelane/pipelaner/pipeline/source"
)

func init() {
	source.RegisterTransform("filter", &Filter{})
}

var (
	ErrInvalidDataType = errors.New("error invalid data type")
)

type EnvMap struct {
	Data any
}

type Filter struct {
	program *vm.Program
}

func (t *Filter) Init(cfg transform.Transform) error {
	filterCfg, ok := cfg.(transform.Filter)
	if !ok {
		return fmt.Errorf("invalid filter config type: %T", cfg)
	}

	program, err := expr.Compile(filterCfg.GetCode(), expr.Env(EnvMap{}))
	if err != nil {
		return err
	}
	t.program = program
	return nil
}

func (t *Filter) Transform(val any) any {
	var v any
	switch value := val.(type) {
	case map[string]any:
		v = value
	case node.AtomicMessage:
		newV := t.Transform(value.Data())
		if err, ok := newV.(error); ok {
			return err
		}
		return value.MessageFrom(newV)
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
	output, err := expr.Run(t.program, EnvMap{Data: v})
	if err != nil {
		return err
	}
	o, ok := output.(bool)
	if !ok || !o {
		return nil
	}
	return val
}
