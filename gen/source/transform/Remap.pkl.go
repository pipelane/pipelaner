// Code generated from Pkl module `pipelaner.source.transforms`. DO NOT EDIT.
package transform

type Remap interface {
	Transform

	GetCode() string
}

var _ Remap = (*RemapImpl)(nil)

type RemapImpl struct {
	SourceName string `pkl:"sourceName"`

	Code string `pkl:"code"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *RemapImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *RemapImpl) GetCode() string {
	return rcv.Code
}

func (rcv *RemapImpl) GetName() string {
	return rcv.Name
}

func (rcv *RemapImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv *RemapImpl) GetThreads() int {
	return rcv.Threads
}

func (rcv *RemapImpl) GetOutputBufferSize() int {
	return rcv.OutputBufferSize
}
