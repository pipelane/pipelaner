package main

import (
	"time"

	"github.com/pipelane/pipelaner/source/generator"
	"github.com/pipelane/pipelaner/source/sink"
	"github.com/pipelane/pipelaner/source/transform"

	"github.com/pipelane/pipelaner"
)

func main() {
	dataSource := pipelane.DataSource{
		Generators: pipelane.Generators{
			"int": &generator.IntGenerator{},
			"cmd": &generator.Command{},
		},
		Maps: pipelane.Maps{
			"inc":  &transform.IncProcessor{},
			"five": &transform.FiveProcessor{},
		},
		Sinks: pipelane.Sinks{
			"console": sink.NewConsole(pipelane.NewLogger()),
		},
	}
	agent, err := pipelane.NewAgent(
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
