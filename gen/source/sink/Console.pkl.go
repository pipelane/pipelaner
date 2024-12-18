// Code generated from Pkl module `com.pipelaner.source.sinks`. DO NOT EDIT.
package sink

import "github.com/pipelane/pipelaner/gen/settings/logger/logformat"

type Console interface {
	Sink

	GetLogFormat() logformat.LogFormat
}

var _ Console = (*ConsoleImpl)(nil)

type ConsoleImpl struct {
	SourceName string `pkl:"sourceName"`

	LogFormat logformat.LogFormat `pkl:"logFormat"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads int `pkl:"threads"`
}

func (rcv *ConsoleImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *ConsoleImpl) GetLogFormat() logformat.LogFormat {
	return rcv.LogFormat
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
