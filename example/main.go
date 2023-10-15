package main

import (
	"time"

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
		dataSource,
		"pipeline.toml",
	)
	if err != nil {
		panic(err)
	}
	go func() {
		time.Sleep(time.Second * 15)
		a.Stop()
	}()
	a.Serve()
}
