/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package utils

import (
	"reflect"
	"sync"

	"github.com/LastPossum/kamino"
)

func MergeChannels(channels []chan any) chan any {
	var wg sync.WaitGroup
	out := make(chan any, len(channels))

	for _, ch := range channels {
		wg.Add(1)
		go func() {
			for msg := range ch {
				out <- msg
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func broadcastChannels(outputs []chan any, ch chan any) {
	if len(outputs) == 0 {
		return
	}
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
		kind := reflect.TypeOf(v).Kind()
		switch kind {
		case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Struct:
			c, err := kamino.Clone(v)
			if err != nil {
				return
			}
			v = c
		default:
		}
		for _, c := range channels {
			c <- v
		}
	}
}
