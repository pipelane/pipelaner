/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package source

import (
	_ "pipelaner/source/generator/cmd"
	_ "pipelaner/source/generator/pipelaner"
	_ "pipelaner/source/sink/console"
	_ "pipelaner/source/sink/pipelaner"
	_ "pipelaner/source/transform/batch"
	_ "pipelaner/source/transform/chunks"
	_ "pipelaner/source/transform/debounce"
	_ "pipelaner/source/transform/filter"
	_ "pipelaner/source/transform/throttling"
)
