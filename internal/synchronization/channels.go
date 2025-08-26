/*
 * Copyright (c) 2023 Alexey Khokhlov
 */

package synchronization

import (
	"reflect"
	"sync"

	"github.com/LastPossum/kamino"
)

func MergeInputs[T any](chs ...chan T) chan T {
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
			for v := range c {
				res <- v
			}
		}(ch)
	}
	return res
}

func BroadcastChannels(outputs []chan any, ch chan any) {
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
