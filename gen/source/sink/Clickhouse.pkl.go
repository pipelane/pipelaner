// Code generated from Pkl module `com.pipelaner.source.sinks`. DO NOT EDIT.
package sink

type Clickhouse interface {
	Sink

	GetAddress() string

	GetUser() string

	GetPassword() string

	GetDatabase() string

	GetTableName() string

	GetAsyncInsert() string

	GetWaitForAsyncInsert() string
}

var _ Clickhouse = (*ClickhouseImpl)(nil)

type ClickhouseImpl struct {
	SourceName string `pkl:"sourceName"`

	Address string `pkl:"address"`

	User string `pkl:"user"`

	Password string `pkl:"password"`

	Database string `pkl:"database"`

	TableName string `pkl:"tableName"`

	AsyncInsert string `pkl:"asyncInsert"`

	WaitForAsyncInsert string `pkl:"waitForAsyncInsert"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads uint `pkl:"threads"`
}

func (rcv *ClickhouseImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv *ClickhouseImpl) GetAddress() string {
	return rcv.Address
}

func (rcv *ClickhouseImpl) GetUser() string {
	return rcv.User
}

func (rcv *ClickhouseImpl) GetPassword() string {
	return rcv.Password
}

func (rcv *ClickhouseImpl) GetDatabase() string {
	return rcv.Database
}

func (rcv *ClickhouseImpl) GetTableName() string {
	return rcv.TableName
}

func (rcv *ClickhouseImpl) GetAsyncInsert() string {
	return rcv.AsyncInsert
}

func (rcv *ClickhouseImpl) GetWaitForAsyncInsert() string {
	return rcv.WaitForAsyncInsert
}

func (rcv *ClickhouseImpl) GetName() string {
	return rcv.Name
}

func (rcv *ClickhouseImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv *ClickhouseImpl) GetThreads() uint {
	return rcv.Threads
}
