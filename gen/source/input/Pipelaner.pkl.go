// Code generated from Pkl module `com.pipelaner.source.inputs`. DO NOT EDIT.
package input

import (
	"github.com/pipelane/pipelaner/gen/source/common"
	"github.com/pipelane/pipelaner/gen/source/input/connectiontype"
)

type Pipelaner interface {
	Input

	GetCommonConfig() *common.Pipelaner

	GetConnectionType() connectiontype.ConnectionType

	GetUnixSocketPath() *string
}

var _ Pipelaner = (*PipelanerImpl)(nil)

type PipelanerImpl struct {
	SourceName string `pkl:"sourceName"`

	CommonConfig *common.Pipelaner `pkl:"commonConfig"`

	ConnectionType connectiontype.ConnectionType `pkl:"connectionType"`

	UnixSocketPath *string `pkl:"unixSocketPath"`

	Name string `pkl:"name"`

	Threads uint `pkl:"threads"`

	OutputBufferSize uint `pkl:"outputBufferSize"`
}

func (rcv *PipelanerImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *PipelanerImpl) GetCommonConfig() *common.Pipelaner {
	return rcv.CommonConfig
}

func (rcv *PipelanerImpl) GetConnectionType() connectiontype.ConnectionType {
	return rcv.ConnectionType
}

func (rcv *PipelanerImpl) GetUnixSocketPath() *string {
	return rcv.UnixSocketPath
}

func (rcv *PipelanerImpl) GetName() string {
	return rcv.Name
}

func (rcv *PipelanerImpl) GetThreads() uint {
	return rcv.Threads
}

func (rcv *PipelanerImpl) GetOutputBufferSize() uint {
	return rcv.OutputBufferSize
}
