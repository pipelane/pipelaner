/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package pipelaner

import (
	"context"
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
