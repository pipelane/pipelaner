// Code generated from Pkl module `pipelaner.source.inputs`. DO NOT EDIT.
package input

type ExampleGenInt interface {
	Input

	GetCount() int
}

var _ ExampleGenInt = (*ExampleGenIntImpl)(nil)

type ExampleGenIntImpl struct {
	SourceName string `pkl:"sourceName"`

	Count int `pkl:"count"`

	Name string `pkl:"name"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *ExampleGenIntImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *ExampleGenIntImpl) GetCount() int {
	return rcv.Count
}

func (rcv *ExampleGenIntImpl) GetName() string {
	return rcv.Name
}

func (rcv *ExampleGenIntImpl) GetThreads() int {
	return rcv.Threads
}

func (rcv *ExampleGenIntImpl) GetOutputBufferSize() int {
	return rcv.OutputBufferSize
}
