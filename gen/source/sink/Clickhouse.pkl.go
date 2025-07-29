// Code generated from Pkl module `com.pipelaner.source.sinks`. DO NOT EDIT.
package sink

import "github.com/pipelane/pipelaner/gen/source/common"

type Clickhouse interface {
	Sink

	GetCredentials() common.ChCredentials

	GetTableName() string

	GetAsyncInsert() string

	GetWaitForAsyncInsert() string

	GetMaxPartitionsPerInsertBlock() int
}

var _ Clickhouse = ClickhouseImpl{}

type ClickhouseImpl struct {
	SourceName string `pkl:"sourceName"`

	Credentials common.ChCredentials `pkl:"credentials"`

	TableName string `pkl:"tableName"`

	AsyncInsert string `pkl:"asyncInsert"`

	WaitForAsyncInsert string `pkl:"waitForAsyncInsert"`

	MaxPartitionsPerInsertBlock int `pkl:"maxPartitionsPerInsertBlock"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads uint `pkl:"threads"`
}

func (rcv ClickhouseImpl) GetSourceName() string {
	return rcv.SourceName
}

func (rcv ClickhouseImpl) GetCredentials() common.ChCredentials {
	return rcv.Credentials
}

func (rcv ClickhouseImpl) GetTableName() string {
	return rcv.TableName
}

func (rcv ClickhouseImpl) GetAsyncInsert() string {
	return rcv.AsyncInsert
}

func (rcv ClickhouseImpl) GetWaitForAsyncInsert() string {
	return rcv.WaitForAsyncInsert
}

func (rcv ClickhouseImpl) GetMaxPartitionsPerInsertBlock() int {
	return rcv.MaxPartitionsPerInsertBlock
}

func (rcv ClickhouseImpl) GetName() string {
	return rcv.Name
}

func (rcv ClickhouseImpl) GetInputs() []string {
	return rcv.Inputs
}

func (rcv ClickhouseImpl) GetThreads() uint {
	return rcv.Threads
}
