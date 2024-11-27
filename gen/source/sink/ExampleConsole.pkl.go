// Code generated from Pkl module `pipelaner.source.sinks`. DO NOT EDIT.
package sink

type ExampleConsole interface {
	Sink
}

var _ ExampleConsole = (*ExampleConsoleImpl)(nil)

type ExampleConsoleImpl struct {
	SourceName string `pkl:"sourceName"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads int `pkl:"threads"`
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

func (rcv *ExampleConsoleImpl) GetThreads() int {
	return rcv.Threads
}
