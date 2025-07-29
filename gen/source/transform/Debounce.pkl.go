// Code generated from Pkl module `com.pipelaner.source.transforms`. DO NOT EDIT.
package transform

import "github.com/apple/pkl-go/pkl"

type Debounce interface {
	Transform

	GetInterval() pkl.Duration
}

var _ Debounce = DebounceImpl{}

type DebounceImpl struct {
	SourceName string `pkl:"sourceName"`

	Interval pkl.Duration `pkl:"interval"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads uint `pkl:"threads"`

	OutputBufferSize uint `pkl:"outputBufferSize"`
}

func (rcv DebounceImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv DebounceImpl) GetInterval() pkl.Duration {
	return rcv.Interval
}

func (rcv DebounceImpl) GetName() string {
	return rcv.Name
}

func (rcv DebounceImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv DebounceImpl) GetThreads() uint {
	return rcv.Threads
}

func (rcv DebounceImpl) GetOutputBufferSize() uint {
	return rcv.OutputBufferSize
}
