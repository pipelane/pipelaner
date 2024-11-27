/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package source

//nolint
import (
	_ "github.com/pipelane/pipelaner/internal/source/generator/cmd"
	_ "github.com/pipelane/pipelaner/internal/source/generator/kafka"
	_ "github.com/pipelane/pipelaner/internal/source/generator/pipelaner"
	_ "github.com/pipelane/pipelaner/internal/source/sink/clickhouse"
	_ "github.com/pipelane/pipelaner/internal/source/sink/console"
	_ "github.com/pipelane/pipelaner/internal/source/sink/kafka"
	_ "github.com/pipelane/pipelaner/internal/source/sink/pipelaner"
	_ "github.com/pipelane/pipelaner/internal/source/transform/batch"
	_ "github.com/pipelane/pipelaner/internal/source/transform/chunks"
	_ "github.com/pipelane/pipelaner/internal/source/transform/debounce"
	_ "github.com/pipelane/pipelaner/internal/source/transform/filter"
	_ "github.com/pipelane/pipelaner/internal/source/transform/remap"
	_ "github.com/pipelane/pipelaner/internal/source/transform/throttling"
)
