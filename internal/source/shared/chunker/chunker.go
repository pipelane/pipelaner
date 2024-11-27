package chunker

import (
	"context"
	"sync/atomic"
	"time"
)

type Config struct {
	MaxChunkSize int
	BufferSize   int
	MaxIdleTime  time.Duration
}
type Chunks[T any] struct {
	Cfg     Config
	buffers chan chan T
	stopped atomic.Bool
	ctx     context.Context
	cancel  context.CancelFunc
	input   chan T
}

func (c *Chunks[T]) Ctx() context.Context {
	return c.ctx
}

func NewChunks[T any](ctx context.Context, cfg Config) *Chunks[T] {
	b := &Chunks[T]{
		Cfg: cfg,
	}
	b.input = make(chan T, b.Cfg.MaxChunkSize*b.Cfg.BufferSize)
	b.resetChannels()
	b.ctx, b.cancel = context.WithCancel(ctx)
	return b
}

func (c *Chunks[T]) Input() chan<- T {
	return c.input
}

func (c *Chunks[T]) Close() {
	if c.stopped.Load() {
		return
	}
	c.stopped.Store(true)
	c.cancel()
	close(c.buffers)
}

func (c *Chunks[T]) resetChannels() {
	c.buffers = make(chan chan T, c.Cfg.BufferSize)
}

func (c *Chunks[T]) send(ch chan T) {
	if !c.stopped.Load() {
		c.buffers <- ch
	}
}
func (c *Chunks[T]) NewChunk() chan T {
	b := make(chan T, c.Cfg.MaxChunkSize)
	c.send(b)
	return b
}

func (c *Chunks[T]) GetChunks() <-chan chan T {
	return c.buffers
}

func (c *Chunks[T]) Generator() {
	counter := atomic.Int64{}
	counter.Store(0)
	timer := time.NewTicker(c.Cfg.MaxIdleTime)
	buffer := c.NewChunk()
	go func() {
		defer c.Close()
		defer timer.Stop()
		for {
			select {
			case <-c.Ctx().Done():
				close(buffer)
				return
			case <-timer.C:
				if counter.Load() == 0 {
					continue
				}
				close(buffer)
				counter.Store(0)
				buffer = c.NewChunk()
			case msg := <-c.input:
				if c.stopped.Load() {
					close(buffer)
					return
				}
				timer.Reset(c.Cfg.MaxIdleTime)
				buffer <- msg
				counter.Add(1)

				if counter.Load() >= int64(c.Cfg.MaxChunkSize) {
					counter.Store(0)
					close(buffer)
					buffer = c.NewChunk()
				}
			}
		}
	}()
}
