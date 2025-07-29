// Code generated from Pkl module `com.pipelaner.source.transforms`. DO NOT EDIT.
package transform

type Batch interface {
	Transform

	GetSize() uint
}

var _ Batch = BatchImpl{}

type BatchImpl struct {
	SourceName string `pkl:"sourceName"`

	Size uint `pkl:"size"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads uint `pkl:"threads"`

	OutputBufferSize uint `pkl:"outputBufferSize"`
}

func (rcv BatchImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv BatchImpl) GetSize() uint {
	return rcv.Size
}

func (rcv BatchImpl) GetName() string {
	return rcv.Name
}

func (rcv BatchImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv BatchImpl) GetThreads() uint {
	return rcv.Threads
}

func (rcv BatchImpl) GetOutputBufferSize() uint {
	return rcv.OutputBufferSize
}
