// Code generated from Pkl module `com.pipelaner.source.transforms`. DO NOT EDIT.
package transform

import "github.com/apple/pkl-go/pkl"

type Throttling interface {
	Transform

	GetInterval() pkl.Duration
}

var _ Throttling = ThrottlingImpl{}

type ThrottlingImpl struct {
	SourceName string `pkl:"sourceName"`

	Interval pkl.Duration `pkl:"interval"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads uint `pkl:"threads"`

	OutputBufferSize uint `pkl:"outputBufferSize"`
}

func (rcv ThrottlingImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv ThrottlingImpl) GetInterval() pkl.Duration {
	return rcv.Interval
}

func (rcv ThrottlingImpl) GetName() string {
	return rcv.Name
}

func (rcv ThrottlingImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv ThrottlingImpl) GetThreads() uint {
	return rcv.Threads
}

func (rcv ThrottlingImpl) GetOutputBufferSize() uint {
	return rcv.OutputBufferSize
}
