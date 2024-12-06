/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package sources

//nolint
import (
	_ "github.com/pipelane/pipelaner/sources/generator/cmd"
	_ "github.com/pipelane/pipelaner/sources/generator/kafka"
	_ "github.com/pipelane/pipelaner/sources/generator/pipelaner"
	_ "github.com/pipelane/pipelaner/sources/sink/clickhouse"
	_ "github.com/pipelane/pipelaner/sources/sink/console"
	_ "github.com/pipelane/pipelaner/sources/sink/kafka"
	_ "github.com/pipelane/pipelaner/sources/sink/pipelaner"
	_ "github.com/pipelane/pipelaner/sources/transform/batch"
	_ "github.com/pipelane/pipelaner/sources/transform/chunks"
	_ "github.com/pipelane/pipelaner/sources/transform/debounce"
	_ "github.com/pipelane/pipelaner/sources/transform/filter"
	_ "github.com/pipelane/pipelaner/sources/transform/remap"
	_ "github.com/pipelane/pipelaner/sources/transform/throttling"
)
