// Code generated from Pkl module `pipelaner.source.examples.custom`. DO NOT EDIT.
package custom

import "github.com/pipelane/pipelaner/gen/source/transform"

type ExampleMul interface {
	transform.Transform

	GetMul() int
}

var _ ExampleMul = (*ExampleMulImpl)(nil)

type ExampleMulImpl struct {
	SourceName string `pkl:"sourceName"`

	Mul int `pkl:"mul"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *ExampleMulImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *ExampleMulImpl) GetMul() int {
	return rcv.Mul
}

func (rcv *ExampleMulImpl) GetName() string {
	return rcv.Name
}

func (rcv *ExampleMulImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv *ExampleMulImpl) GetThreads() int {
	return rcv.Threads
}

func (rcv *ExampleMulImpl) GetOutputBufferSize() int {
	return rcv.OutputBufferSize
}
