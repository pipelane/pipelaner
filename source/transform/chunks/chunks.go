/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package chunks

/*import (
	"sync/atomic"
	"time"

	"github.com/pipelane/pipelaner"
	"github.com/pipelane/pipelaner/source/shared/chunker"
)

type ChunkCfg struct {
	MaxChunkSize int64  `pipelane:"max_chunk_size"`
	MaxIdleTime  string `pipelane:"max_idle_time"`
}

func (c *ChunkCfg) Interval() (time.Duration, error) {
	return time.ParseDuration(c.MaxIdleTime)
}

type Chunk struct {
	cfg    *pipelaner.BaseLaneConfig
	buffer *chunker.Chunks[any]
	locked atomic.Bool
}

func init() {
	pipelaner.RegisterMap("chunks", &Chunk{})
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
	c.buffer = chunker.NewChunks[any](ctx.Context(), chunker.Config{
		MaxChunkSize: int(v.MaxChunkSize),
		BufferSize:   int(c.cfg.OutputBufferSize),
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
}*/
