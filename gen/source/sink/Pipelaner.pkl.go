// Code generated from Pkl module `pipelaner.source.sinks`. DO NOT EDIT.
package sink

type Pipelaner interface {
	Sink

	GetHost() string

	GetPort() int

	GetTls() bool

	GetCertFile() *string

	GetKeyFile() *string
}

var _ Pipelaner = (*PipelanerImpl)(nil)

type PipelanerImpl struct {
	SourceName string `pkl:"sourceName"`

	Host string `pkl:"host"`

	Port int `pkl:"port"`

	Tls bool `pkl:"tls"`

	CertFile *string `pkl:"certFile"`

	KeyFile *string `pkl:"keyFile"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads int `pkl:"threads"`
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

func (rcv *PipelanerImpl) GetTls() bool {
	return rcv.Tls
}

func (rcv *PipelanerImpl) GetCertFile() *string {
	return rcv.CertFile
}

func (rcv *PipelanerImpl) GetKeyFile() *string {
	return rcv.KeyFile
}

func (rcv *PipelanerImpl) GetName() string {
	return rcv.Name
}

func (rcv *PipelanerImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv *PipelanerImpl) GetThreads() int {
	return rcv.Threads
}
