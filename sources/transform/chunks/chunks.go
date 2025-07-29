/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package chunks

import (
	"fmt"
	"sync/atomic"

	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/pipeline/source"
	"github.com/pipelane/pipelaner/sources/shared/chunker"
)

func init() {
	source.RegisterTransform("chunks", &Chunk{})
}

type Chunk struct {
	buffer *chunker.Chunks
	locked atomic.Bool
}

func (c *Chunk) Init(cfg transform.Transform) error {
	chunkCfg, ok := cfg.(transform.Chunk)
	if !ok {
		return fmt.Errorf("invalid chunk config type: %T", cfg)
	}

	time := chunkCfg.GetMaxIdleTime()
	c.buffer = chunker.NewChunks(chunker.Config{
		MaxChunkSize: chunkCfg.GetMaxChunkSize(),
		BufferSize:   chunkCfg.GetOutputBufferSize(),
		MaxIdleTime:  time.GoDuration(),
	})
	c.buffer.Generator()
	return nil
}

func (c *Chunk) Transform(val any) any {
	c.buffer.SetValue(val)
	if c.locked.Load() {
		return nil
	}
	c.locked.Store(true)
	defer c.locked.Store(false)
	return c.buffer.Chunk()
}
