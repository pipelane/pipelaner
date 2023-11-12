/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package main

import (
	"time"

	"pipelaner/source/generator"
	"pipelaner/source/sink"
	"pipelaner/source/transform"

	"github.com/pipelane/pipelaner"
)

func main() {
	dataSource := pipelaner.DataSource{
		Generators: pipelaner.Generators{
			"exec": &generator.Exec{},
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
			"console": sink.NewConsole(pipelaner.NewLogger()),
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
