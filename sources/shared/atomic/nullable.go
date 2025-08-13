/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package atomic

import (
	"errors"
	"sync/atomic"
)

type atomicNullableStoredValue struct {
	val any
}

type Nullable struct {
	val atomic.Value
}

func (a *Nullable) Store(v any) {
	a.val.Store(atomicNullableStoredValue{val: v})
}

func (a *Nullable) Load() any {
	v, ok := a.val.Load().(atomicNullableStoredValue)
	if !ok {
		panic(errors.New("atomic nullable stored value not initialized"))
	}
	return v.val
}
