/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package synchronization

type Semaphore struct {
	c chan struct{}
}

func NewSemaphore(n uint) *Semaphore {
	return &Semaphore{make(chan struct{}, n)}
}

func (s *Semaphore) Acquire() {
	s.c <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.c
}
