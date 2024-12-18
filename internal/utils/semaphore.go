/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package utils

type Semaphore struct {
	c chan struct{}
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{make(chan struct{}, n)}
}

func (s *Semaphore) Acquire() {
	s.c <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.c
}
