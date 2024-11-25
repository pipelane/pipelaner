// Code generated from Pkl module `pipelaner.source.inputs`. DO NOT EDIT.
package input

type GenInt interface {
	Input

	GetCount() int
}

var _ GenInt = (*GenIntImpl)(nil)

type GenIntImpl struct {
	SourceName string `pkl:"sourceName"`

	Count int `pkl:"count"`

	Name string `pkl:"name"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *GenIntImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *GenIntImpl) GetCount() int {
	return rcv.Count
}

func (rcv *GenIntImpl) GetName() string {
	return rcv.Name
}

func (rcv *GenIntImpl) GetThreads() int {
	return rcv.Threads
}

func (rcv *GenIntImpl) GetOutputBufferSize() int {
	return rcv.OutputBufferSize
}
