// Code generated from Pkl module `pipelaner.source.examples.custom`. DO NOT EDIT.
package custom

import "github.com/pipelane/pipelaner/gen/source/sink"

type ExampleConsole interface {
	sink.Sink
}

var _ ExampleConsole = (*ExampleConsoleImpl)(nil)

type ExampleConsoleImpl struct {
	SourceName string `pkl:"sourceName"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads uint `pkl:"threads"`
}

func (rcv *ExampleConsoleImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *ExampleConsoleImpl) GetName() string {
	return rcv.Name
}

func (rcv *ExampleConsoleImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv *ExampleConsoleImpl) GetThreads() uint {
	return rcv.Threads
}
