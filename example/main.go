/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package main

import (
	"time"

	"pipelaner"
)

func main() {
	//dataSource := pipelaner.DataSource{
	//	Generators: pipelaner.Generators{
	//		"exec":      &generator.Exec{},
	//		"pipelaner": &generator.Pipelaner{},
	//		// For Test
	//		"int":  &tests.IntGenerator{},
	//		"int2": &tests.IntTwoGenerator{},
	//		"rand": &tests.MapGenerator{},
	//	},
	//	Maps: pipelaner.Maps{
	//		"filter":     &transform.Filter{},
	//		"debounce":   &transform.Debounce{},
	//		"throttling": &transform.Throttling{},
	//		"batch":      &transform.Batch{},
	//		"chunks":     &transform.Chunk{},
	//		// For Test
	//		"int_tr":   &tests.IntTransform{},
	//		"int_tr_e": &tests.IntTransformEmpty{},
	//	},
	//	Sinks: pipelaner.Sinks{
	//		"console":   sink.NewConsole(pipelaner.NewLogger()),
	//		"pipelaner": &sink.Pipelaner{},
	//	},
	//}
	agent, err := pipelaner.NewAgent(
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
