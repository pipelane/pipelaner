// Code generated from Pkl module `com.pipelaner.source.sinks`. DO NOT EDIT.
package sink

type Console interface {
	Sink
}

var _ Console = (*ConsoleImpl)(nil)

type ConsoleImpl struct {
	SourceName string `pkl:"sourceName"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads int `pkl:"threads"`
}

func (rcv *ConsoleImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *ConsoleImpl) GetName() string {
	return rcv.Name
}

func (rcv *ConsoleImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv *ConsoleImpl) GetThreads() int {
	return rcv.Threads
}
