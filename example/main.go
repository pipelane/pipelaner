/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package main

import (
	"pipelaner/example/tests"
	"time"

	"pipelaner"
	_ "pipelaner/source"
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
