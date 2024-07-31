/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package source

import (
	_ "github.com/pipelane/pipelaner/source/generator/cmd"
	_ "github.com/pipelane/pipelaner/source/generator/kafka"
	_ "github.com/pipelane/pipelaner/source/generator/pipelaner"
	_ "github.com/pipelane/pipelaner/source/sink/clickhouse"
	_ "github.com/pipelane/pipelaner/source/sink/console"
	_ "github.com/pipelane/pipelaner/source/sink/kafka"
	_ "github.com/pipelane/pipelaner/source/sink/pipelaner"
	_ "github.com/pipelane/pipelaner/source/transform/batch"
	_ "github.com/pipelane/pipelaner/source/transform/chunks"
	_ "github.com/pipelane/pipelaner/source/transform/debounce"
	_ "github.com/pipelane/pipelaner/source/transform/filter"
	_ "github.com/pipelane/pipelaner/source/transform/throttling"
)
