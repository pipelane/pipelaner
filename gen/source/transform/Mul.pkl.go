// Code generated from Pkl module `pipelaner.source.transforms`. DO NOT EDIT.
package transform

type Mul interface {
	Transform

	GetMul() int
}

var _ Mul = (*MulImpl)(nil)

type MulImpl struct {
	SourceName string `pkl:"sourceName"`

	Mul int `pkl:"mul"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *MulImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *MulImpl) GetMul() int {
	return rcv.Mul
}

func (rcv *MulImpl) GetName() string {
	return rcv.Name
}

func (rcv *MulImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv *MulImpl) GetThreads() int {
	return rcv.Threads
}

func (rcv *MulImpl) GetOutputBufferSize() int {
	return rcv.OutputBufferSize
}
