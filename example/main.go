/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package main

import (
	"log"

	"github.com/pipelane/pipelaner"
	"github.com/pipelane/pipelaner/example/tests"
	_ "github.com/pipelane/pipelaner/source"
)

func main() {
	log.Println("Start")
	pipelaner.RegisterGenerator("int", &tests.IntGenerator{})
	agent, err := pipelaner.NewAgent(
		"example/pipeline.toml",
	)
	if err != nil {
		panic(err)
	}
	agent.Serve()
}
