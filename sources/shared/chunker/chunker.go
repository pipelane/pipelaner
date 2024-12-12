package chunker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	MaxChunkSize int
	BufferSize   int
	MaxIdleTime  time.Duration
}
type Chunks struct {
	Cfg     Config
	buffers chan chan any
	input   chan any
	wg      sync.WaitGroup
	cancel  context.CancelFunc
}

func NewChunks(cfg Config) *Chunks {
	b := &Chunks{
		Cfg:   cfg,
		input: make(chan any, cfg.MaxChunkSize*cfg.BufferSize),
		wg:    sync.WaitGroup{},
	}
	b.resetChannels()
	return b
}

func (c *Chunks) SetValue(v any) {
	c.wg.Add(1)
	c.input <- v
}

func (c *Chunks) resetChannels() {
	c.buffers = make(chan chan any, c.Cfg.BufferSize)
}

func (c *Chunks) NewChunk() chan any {
	b := make(chan any, c.Cfg.MaxChunkSize)
	c.buffers <- b
	return b
}

func (c *Chunks) Chunks() <-chan chan any {
	return c.buffers
}
func (c *Chunks) Chunk() chan any {
	return <-c.buffers
}
func (c *Chunks) Stop() {
	c.cancel()
}

func (c *Chunks) Generator() {
	timer := time.NewTimer(c.Cfg.MaxIdleTime)
	stop := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	go c.startProcessing(stop, timer)
	go func() {
		<-ctx.Done()
		c.wg.Wait()
		close(c.input)
		timer.Stop()
		close(c.buffers)
		stop <- struct{}{}
		close(stop)
	}()
}

func (c *Chunks) startProcessing(onClose chan struct{}, timer *time.Timer) {
	buffer := c.NewChunk()
	counter := atomic.Int64{}
	counter.Store(0)
Loop:
	for {
		select {
		case <-onClose:
			break Loop
		case <-timer.C:
			if counter.Load() == 0 {
				timer.Reset(c.Cfg.MaxIdleTime)
				continue
			}
			close(buffer)
			counter.Store(0)
			buffer = c.NewChunk()
			timer.Reset(c.Cfg.MaxIdleTime)
		case msg, ok := <-c.input:
			if !ok && msg == nil {
				continue
			}
			timer.Reset(c.Cfg.MaxIdleTime)
			buffer <- msg
			counter.Add(1)
			if counter.Load() == int64(c.Cfg.MaxChunkSize) {
				close(buffer)
				counter.Store(0)
				buffer = c.NewChunk()
			}
			c.wg.Done()
		}
	}
	close(buffer)
}
