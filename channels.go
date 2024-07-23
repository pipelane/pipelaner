/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
	"github.com/LastPossum/kamino"
	"reflect"
	"sync"
)

func mergeInputs[T any](ctx context.Context, chs ...chan T) chan T {
	if len(chs) == 1 {
		return chs[0]
	}
	lens := 0
	for i := range chs {
		lens += cap(chs[i])
	}
	res := make(chan T, lens)
	gr := sync.WaitGroup{}
	gr.Add(len(chs))
	go func() {
		gr.Wait()
		close(res)
	}()
	for _, ch := range chs {
		go func(c chan T) {
			defer gr.Done()
			for {
				select {
				case <-ctx.Done():
					break
				case v, ok := <-c:
					if !ok {
						break
					}
					res <- v
				}
			}
		}(ch)
	}
	return res
}

func broadcastChannels(outputs []chan any, ch chan any) {
	channels := make([]chan any, len(outputs))
	for i := 0; i < len(channels); i++ {
		channels[i] = make(chan any, cap(ch))
	}
	defer func() {
		for _, c := range channels {
			close(c)
		}
	}()

	for i := range outputs {
		outputs[i] <- channels[i]
	}

	for v := range ch {
		val := reflect.TypeOf(v)
		if val.Kind() == reflect.Pointer {
			c, err := kamino.Clone(v)
			if err != nil {
				return
			}
			v = c
		}
		for _, c := range channels {
			c <- v
		}
	}
}
