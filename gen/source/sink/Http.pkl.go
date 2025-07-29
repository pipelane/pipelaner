// Code generated from Pkl module `com.pipelaner.source.sinks`. DO NOT EDIT.
package sink

import "github.com/pipelane/pipelaner/gen/source/sink/method"

type Http interface {
	Sink

	GetUrl() string

	GetMethod() method.Method

	GetHeaders() *map[string]string
}

var _ Http = HttpImpl{}

type HttpImpl struct {
	SourceName string `pkl:"sourceName"`

	Url string `pkl:"url"`

	Method method.Method `pkl:"method"`

	Headers *map[string]string `pkl:"headers"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads uint `pkl:"threads"`
}

func (rcv HttpImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv HttpImpl) GetUrl() string {
	return rcv.Url
}

func (rcv HttpImpl) GetMethod() method.Method {
	return rcv.Method
}

func (rcv HttpImpl) GetHeaders() *map[string]string {
	return rcv.Headers
}

func (rcv HttpImpl) GetName() string {
	return rcv.Name
}

func (rcv HttpImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv HttpImpl) GetThreads() uint {
	return rcv.Threads
}
