/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package main

import (
	"time"

	"github.com/pipelane/pipelaner"
	"github.com/pipelane/pipelaner/source/generator"
	"github.com/pipelane/pipelaner/source/sink"
	"github.com/pipelane/pipelaner/source/transform"
)

func main() {
	dataSource := pipelaner.DataSource{
		Generators: pipelaner.Generators{
			"exec":      &generator.Exec{},
			"pipelaner": &generator.Pipelaner{},
			//For Test
			"int":  &generator.IntGenerator{},
			"rand": &generator.MapGenerator{},
		},
		Maps: pipelaner.Maps{
			"filter": &transform.Filter{},
			//For Test
			"inc":  &transform.IncProcessor{},
			"five": &transform.FiveProcessor{},
		},
		Sinks: pipelaner.Sinks{
			"console":   sink.NewConsole(pipelaner.NewLogger()),
			"pipelaner": &sink.Pipelaner{},
		},
	}
	agent, err := pipelaner.NewAgent(
		dataSource,
		"pipeline.toml",
	)
	if err != nil {
		panic(err)
	}
	go func() {
		time.Sleep(time.Second * 15)
		agent.Stop()
	}()
	agent.Serve()
}
