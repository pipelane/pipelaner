/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package main

import (
	"context"
	"time"

	"github.com/pipelane/pipelaner"
	_ "github.com/pipelane/pipelaner/sources"
)

func main() {
	ctx := context.Background()
	agent, err := pipelaner.NewAgent(
		"examples/custom/pkl/config.pkl",
	)
	if err != nil {
		panic(err)
	}
	lock := make(chan struct{})
	go func() {
		time.Sleep(time.Second * 15)
		err = agent.Shutdown(context.Background())
		if err != nil {
			panic(err)
		}
		lock <- struct{}{}
	}()
	go func() {
		if err = agent.Serve(ctx); err != nil {
			panic(err)
		}
	}()
	<-lock
}
