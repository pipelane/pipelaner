/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package chunks

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/pipelane/pipelaner"
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

func init() {
	pipelaner.RegisterMap("chunks", &Chunk{})
}

func (c *Chunks[T]) Ctx() context.Context {
	return c.ctx
}

func NewChunks[T any](ctx context.Context, cfg Config) *Chunks[T] {
	b := &Chunks[T]{
		Cfg: cfg,
	}
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
	c.input = make(chan T, c.Cfg.MaxChunkSize*c.Cfg.BufferSize)
	go func() {
		defer c.Close()
		defer timer.Stop()
		for {
			select {
			case <-c.Ctx().Done():
				close(buffer)
				return
			case <-timer.C:
				close(buffer)
				counter.Store(0)
				buffer = c.NewChunk()
			case msg := <-c.input:
				if c.stopped.Load() {
					close(buffer)
					return
				}
				timer.Reset(c.Cfg.MaxIdleTime)
				if counter.Load() >= int64(c.Cfg.MaxChunkSize) {
					counter.Store(0)
					close(buffer)
					buffer = c.NewChunk()
				}
				buffer <- msg
				counter.Add(1)
			}
		}
	}()
}

type ChunkCfg struct {
	MaxChunkSize int64  `pipelane:"max_chunk_size"`
	MaxIdleTime  string `pipelane:"max_idle_time"`
}

func (c *ChunkCfg) Interval() (time.Duration, error) {
	return time.ParseDuration(c.MaxIdleTime)
}

type Chunk struct {
	cfg    *pipelaner.BaseLaneConfig
	buffer *Chunks[any]
	locked atomic.Bool
}

func (c *Chunk) Init(ctx *pipelaner.Context) error {
	c.cfg = ctx.LaneItem().Config()
	v := &ChunkCfg{}
	err := c.cfg.ParseExtended(v)
	if err != nil {
		return err
	}
	interval, err := v.Interval()
	if err != nil {
		return err
	}
	c.buffer = NewChunks[any](ctx.Context(), Config{
		MaxChunkSize: int(v.MaxChunkSize),
		BufferSize:   int(c.cfg.BufferSize),
		MaxIdleTime:  interval,
	})
	c.buffer.Generator()
	return nil
}

func (c *Chunk) Map(_ *pipelaner.Context, val any) any {
	c.buffer.Input() <- val
	if c.locked.Load() {
		return nil
	}
	c.locked.Store(true)
	defer c.locked.Store(false)
	v := <-c.buffer.GetChunks()
	return v
}
