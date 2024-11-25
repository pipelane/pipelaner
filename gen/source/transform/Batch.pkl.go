// Code generated from Pkl module `pipelaner.source.transforms`. DO NOT EDIT.
package transform

type Batch interface {
	Transform

	GetSize() uint32
}

var _ Batch = (*BatchImpl)(nil)

type BatchImpl struct {
	SourceName string `pkl:"sourceName"`

	Size uint32 `pkl:"size"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *BatchImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *BatchImpl) GetSize() uint32 {
	return rcv.Size
}

func (rcv *BatchImpl) GetName() string {
	return rcv.Name
}

func (rcv *BatchImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv *BatchImpl) GetThreads() int {
	return rcv.Threads
}

func (rcv *BatchImpl) GetOutputBufferSize() int {
	return rcv.OutputBufferSize
}
