/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package chunks

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/pipelane/pipelaner/gen/source/transform"
	"github.com/pipelane/pipelaner/internal/pipeline/source"
	"github.com/pipelane/pipelaner/source/shared/chunker"
)

func init() {
	source.RegisterTransform("chunks", &Chunk{})
}

type Chunk struct {
	buffer *chunker.Chunks[any]
	locked atomic.Bool
}

func (c *Chunk) Init(cfg transform.Transform) error {
	chunkCfg, ok := cfg.(transform.Chunk)
	if !ok {
		return fmt.Errorf("invalid chunk config type: %T", cfg)
	}

	c.buffer = chunker.NewChunks[any](context.Background(), chunker.Config{
		MaxChunkSize: int(chunkCfg.GetMaxChunkSize()),
		BufferSize:   chunkCfg.GetOutputBufferSize(),
		MaxIdleTime:  chunkCfg.GetMaxIdleTime().GoDuration(),
	})
	c.buffer.Generator()
	return nil
}

func (c *Chunk) Transform(val any) any {
	c.buffer.Input() <- val
	if c.locked.Load() {
		return nil
	}
	c.locked.Store(true)
	defer c.locked.Store(false)
	v := <-c.buffer.GetChunks()
	return v
}
