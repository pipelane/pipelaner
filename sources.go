/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package pipelaner

var dataSource = newDataSources()

func RegisterSink(name string, sink Sink) {
	dataSource.Sinks[name] = sink
}

func RegisterMap(name string, maps Map) {
	dataSource.Maps[name] = maps
}

func RegisterGenerator(name string, generators Generator) {
	dataSource.Generators[name] = generators
}
