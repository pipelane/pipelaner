// Code generated from Pkl module `com.pipelaner.source.transforms`. DO NOT EDIT.
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

	Threads uint `pkl:"threads"`

	OutputBufferSize uint `pkl:"outputBufferSize"`
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

func (rcv *FilterImpl) GetThreads() uint {
	return rcv.Threads
}

func (rcv *FilterImpl) GetOutputBufferSize() uint {
	return rcv.OutputBufferSize
}
