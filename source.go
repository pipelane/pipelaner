/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

type Generators map[string]Generator
type Maps map[string]Map
type Sinks map[string]Sink

type dataSources struct {
	Generators Generators
	Maps       Maps
	Sinks      Sinks
}

func newDataSources() *dataSources {
	return &dataSources{
		Generators: map[string]Generator{},
		Maps:       map[string]Map{},
		Sinks:      map[string]Sink{},
	}
}
