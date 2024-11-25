// Code generated from Pkl module `pipelaner.source.inputs`. DO NOT EDIT.
package input

type Cmd interface {
	Input

	GetExec() []string
}

var _ Cmd = (*CmdImpl)(nil)

type CmdImpl struct {
	SourceName string `pkl:"sourceName"`

	Exec []string `pkl:"exec"`

	Name string `pkl:"name"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *CmdImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *CmdImpl) GetExec() []string {
	return rcv.Exec
}

func (rcv *CmdImpl) GetName() string {
	return rcv.Name
}

func (rcv *CmdImpl) GetThreads() int {
	return rcv.Threads
}

func (rcv *CmdImpl) GetOutputBufferSize() int {
	return rcv.OutputBufferSize
}
