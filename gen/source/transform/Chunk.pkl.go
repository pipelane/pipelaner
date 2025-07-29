// Code generated from Pkl module `com.pipelaner.source.transforms`. DO NOT EDIT.
package transform

import "github.com/apple/pkl-go/pkl"

type Chunk interface {
	Transform

	GetMaxChunkSize() uint

	GetMaxIdleTime() pkl.Duration
}

var _ Chunk = ChunkImpl{}

type ChunkImpl struct {
	SourceName string `pkl:"sourceName"`

	MaxChunkSize uint `pkl:"maxChunkSize"`

	MaxIdleTime pkl.Duration `pkl:"maxIdleTime"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads uint `pkl:"threads"`

	OutputBufferSize uint `pkl:"outputBufferSize"`
}

func (rcv ChunkImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv ChunkImpl) GetMaxChunkSize() uint {
	return rcv.MaxChunkSize
}

func (rcv ChunkImpl) GetMaxIdleTime() pkl.Duration {
	return rcv.MaxIdleTime
}

func (rcv ChunkImpl) GetName() string {
	return rcv.Name
}

func (rcv ChunkImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv ChunkImpl) GetThreads() uint {
	return rcv.Threads
}

func (rcv ChunkImpl) GetOutputBufferSize() uint {
	return rcv.OutputBufferSize
}
