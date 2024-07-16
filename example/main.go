/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package main

import (
	"github.com/pipelane/pipelaner/example/tests"
	"time"

	"github.com/pipelane/pipelaner"
	_ "github.com/pipelane/pipelaner/source"
)

func main() {
	pipelaner.RegisterGenerator("int", &tests.IntGenerator{})
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
