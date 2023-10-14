package main

import (
	"context"

	"pipelaner/source/generator"
	"pipelaner/source/sink"
	"pipelaner/source/transform"

	"pipelaner"
)

func main() {
	dataSource := pipelane.DataSource{
		Generators: pipelane.Generators{
			"int": &generator.IntGenerator{},
		},
		Maps: pipelane.Maps{
			"inc":  &transform.IncProcessor{},
			"five": &transform.FiveProcessor{},
		},
		Sinks: pipelane.Sinks{
			"console": sink.NewConsole(pipelane.NewLogger()),
		},
	}
	a, err := pipelane.NewAgent(
		context.Background(),
		dataSource,
		"pipeline.toml",
	)
	if err != nil {
		panic(err)
	}
	a.Serve()
}
