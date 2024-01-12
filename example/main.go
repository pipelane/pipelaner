/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package main

import (
	"time"

	"pipelaner"
	"pipelaner/example/tests"
	"pipelaner/source/generator"
	"pipelaner/source/sink"
	"pipelaner/source/transform"
)

func main() {
	dataSource := pipelaner.DataSource{
		Generators: pipelaner.Generators{
			"exec":      &generator.Exec{},
			"pipelaner": &generator.Pipelaner{},
			// For Test
			"int":  &tests.IntGenerator{},
			"int2": &tests.IntTwoGenerator{},
			"rand": &tests.MapGenerator{},
		},
		Maps: pipelaner.Maps{
			"filter":     &transform.Filter{},
			"debounce":   &transform.Debounce{},
			"throttling": &transform.Throttling{},
			"batch":      &transform.Batch{},
			"chunks":     &transform.Chunk{},
			// For Test
			"int_tr":   &tests.IntTransform{},
			"int_tr_e": &tests.IntTransformEmpty{},
		},
		Sinks: pipelaner.Sinks{
			"console":   sink.NewConsole(pipelaner.NewLogger()),
			"pipelaner": &sink.Pipelaner{},
		},
	}
	agent, err := pipelaner.NewAgent(
		dataSource,
		"example/pipeline.toml",
	)
	if err != nil {
		panic(err)
	}
	go func() {
		time.Sleep(time.Second * 120)
		agent.Stop()
	}()
	agent.Serve()
}
