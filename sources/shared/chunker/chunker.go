package chunker

import (
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
	stopped atomic.Bool
	input   chan any
	wg      sync.WaitGroup
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

func (c *Chunks) send(ch chan any) {
	c.buffers <- ch
}
func (c *Chunks) NewChunk() chan any {
	b := make(chan any, c.Cfg.MaxChunkSize)
	c.send(b)
	return b
}

func (c *Chunks) Chunks() <-chan chan any {
	return c.buffers
}
func (c *Chunks) Chunk() chan any {
	return <-c.buffers
}
func (c *Chunks) Stop() {
	c.stopped.Store(true)
}

func (c *Chunks) Generator() {
	timer := time.NewTicker(c.Cfg.MaxIdleTime)
	stop := make(chan struct{}, 1)
	go func() {
		for {
			if c.stopped.Load() {
				c.wg.Wait()
				timer.Stop()
				close(c.buffers)
				close(c.input)
				stop <- struct{}{}
				close(stop)
				break
			}
		}
	}()
	go c.startProcessing(stop, timer)
}

func (c *Chunks) startProcessing(stop chan struct{}, timer *time.Ticker) {
	buffer := c.NewChunk()
	counter := atomic.Int64{}
	counter.Store(0)
	breaks := false
	for !breaks {
		select {
		case <-stop:
			close(buffer)
			breaks = true
			continue
		case <-timer.C:
			if counter.Load() == 0 {
				timer.Reset(c.Cfg.MaxIdleTime)
				continue
			}
			close(buffer)
			counter.Store(0)
			buffer = c.NewChunk()
			timer.Reset(c.Cfg.MaxIdleTime)
		default:
			msg, ok := <-c.input
			if !ok && msg == nil {
				continue
			}
			timer.Reset(c.Cfg.MaxIdleTime)
			buffer <- msg
			counter.Add(1)
			if counter.Load() >= int64(c.Cfg.MaxChunkSize) {
				counter.Store(0)
				close(buffer)
				buffer = c.NewChunk()
			}
			c.wg.Done()
		}
	}
}
