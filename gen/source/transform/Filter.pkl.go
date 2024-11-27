// Code generated from Pkl module `pipelaner.source.transforms`. DO NOT EDIT.
package transform

type Filter interface {
	Transform

	GetCode() string
}

var _ Filter = (*FilterImpl)(nil)

type FilterImpl struct {
	SourceName string `pkl:"sourceName"`

	Code string `pkl:"code"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *FilterImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *FilterImpl) GetCode() string {
	return rcv.Code
}

func (rcv *FilterImpl) GetName() string {
	return rcv.Name
}

func (rcv *FilterImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv *FilterImpl) GetThreads() int {
	return rcv.Threads
}

func (rcv *FilterImpl) GetOutputBufferSize() int {
	return rcv.OutputBufferSize
}