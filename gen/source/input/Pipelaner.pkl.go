// Code generated from Pkl module `pipelaner.source.inputs`. DO NOT EDIT.
package input

type Pipelaner interface {
	Input

	GetHost() *string

	GetPort() int

	GetTls() *bool

	GetCertFile() *string

	GetKeyFile() *string

	GetConnectionType() *string

	GetUnixSocketPasth() string
}

var _ Pipelaner = (*PipelanerImpl)(nil)

type PipelanerImpl struct {
	SourceName string `pkl:"sourceName"`

	Host *string `pkl:"host"`

	Port int `pkl:"port"`

	Tls *bool `pkl:"tls"`

	CertFile *string `pkl:"certFile"`

	KeyFile *string `pkl:"keyFile"`

	ConnectionType *string `pkl:"connectionType"`

	UnixSocketPasth string `pkl:"unixSocketPasth"`

	Name string `pkl:"name"`

	Threads int `pkl:"threads"`

	OutputBufferSize int `pkl:"outputBufferSize"`
}

func (rcv *PipelanerImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *PipelanerImpl) GetHost() *string {
	return rcv.Host
}

func (rcv *PipelanerImpl) GetPort() int {
	return rcv.Port
}

func (rcv *PipelanerImpl) GetTls() *bool {
	return rcv.Tls
}

func (rcv *PipelanerImpl) GetCertFile() *string {
	return rcv.CertFile
}

func (rcv *PipelanerImpl) GetKeyFile() *string {
	return rcv.KeyFile
}

func (rcv *PipelanerImpl) GetConnectionType() *string {
	return rcv.ConnectionType
}

func (rcv *PipelanerImpl) GetUnixSocketPasth() string {
	return rcv.UnixSocketPasth
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
