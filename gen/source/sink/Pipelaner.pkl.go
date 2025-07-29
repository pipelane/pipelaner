// Code generated from Pkl module `com.pipelaner.source.sinks`. DO NOT EDIT.
package sink

import "github.com/pipelane/pipelaner/gen/source/common"

type Pipelaner interface {
	Sink

	GetCommonConfig() common.Pipelaner
}

var _ Pipelaner = PipelanerImpl{}

type PipelanerImpl struct {
	SourceName string `pkl:"sourceName"`

	CommonConfig common.Pipelaner `pkl:"commonConfig"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads uint `pkl:"threads"`
}

func (rcv PipelanerImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv PipelanerImpl) GetCommonConfig() common.Pipelaner {
	return rcv.CommonConfig
}

func (rcv PipelanerImpl) GetName() string {
	return rcv.Name
}

func (rcv PipelanerImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv PipelanerImpl) GetThreads() uint {
	return rcv.Threads
}
