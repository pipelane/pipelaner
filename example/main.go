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
		Transforms: pipelane.Transformers{
			"inc":  &transform.IncProcessor{},
			"five": &transform.FiveProcessor{},
		},
		Sinks: pipelane.Sinks{
			"console": sink.NewConsole(pipelane.NewLogger()),
		},
	}
	ch := make(chan bool, 1)
	_, err := pipelane.NewTreeFrom(
		context.Background(),
		dataSource,
		"pipeline.toml",
	)
	if err != nil {
		panic(err)
	}
	<-ch
}
