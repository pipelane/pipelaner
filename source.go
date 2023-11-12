/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

type Generators map[string]Generator
type Maps map[string]Map
type Sinks map[string]Sink

type DataSource struct {
	Generators Generators
	Maps       Maps
	Sinks      Sinks
}
