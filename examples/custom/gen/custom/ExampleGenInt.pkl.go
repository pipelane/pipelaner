// Code generated from Pkl module `pipelaner.source.examples.custom`. DO NOT EDIT.
package custom

import "github.com/pipelane/pipelaner/gen/source/input"

type ExampleGenInt interface {
	input.Input

	GetCount() int
}

var _ ExampleGenInt = (*ExampleGenIntImpl)(nil)

type ExampleGenIntImpl struct {
	SourceName string `pkl:"sourceName"`

	Count int `pkl:"count"`

	Name string `pkl:"name"`

	Threads uint `pkl:"threads"`

	OutputBufferSize uint `pkl:"outputBufferSize"`
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

func (rcv *ExampleGenIntImpl) GetThreads() uint {
	return rcv.Threads
}

func (rcv *ExampleGenIntImpl) GetOutputBufferSize() uint {
	return rcv.OutputBufferSize
}
