// Code generated from Pkl module `pipelaner.source.inputs`. DO NOT EDIT.
package input

import "github.com/pipelane/pipelaner/gen/source/input/connectiontype"

type Pipelaner interface {
	Input

	GetHost() string

	GetPort() int

	GetTls() *TLSConfig

	GetConnectionType() connectiontype.ConnectionType

	GetUnixSocketPath() *string
}

var _ Pipelaner = (*PipelanerImpl)(nil)

type PipelanerImpl struct {
	SourceName string `pkl:"sourceName"`

	Host string `pkl:"host"`

	Port int `pkl:"port"`

	Tls *TLSConfig `pkl:"tls"`

	ConnectionType connectiontype.ConnectionType `pkl:"connectionType"`

	UnixSocketPath *string `pkl:"unixSocketPath"`

	Name string `pkl:"name"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *PipelanerImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *PipelanerImpl) GetHost() string {
	return rcv.Host
}

func (rcv *PipelanerImpl) GetPort() int {
	return rcv.Port
}

func (rcv *PipelanerImpl) GetTls() *TLSConfig {
	return rcv.Tls
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

func (rcv *PipelanerImpl) GetThreads() int {
	return rcv.Threads
}

func (rcv *PipelanerImpl) GetOutputBufferSize() int {
	return rcv.OutputBufferSize
}
