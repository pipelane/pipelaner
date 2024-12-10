// Code generated from Pkl module `com.pipelaner.source.sinks`. DO NOT EDIT.
package sink

import "github.com/apple/pkl-go/pkl"

type Clickhouse interface {
	Sink

	GetAddress() string

	GetUser() string

	GetPassword() string

	GetDatabase() string

	GetMigrationEngine() string

	GetMigrationsPathClickhouse() string

	GetMaxExecutionTime() *pkl.Duration

	GetCannMaxLifeTime() *pkl.Duration

	GetDialTimeout() *pkl.Duration

	GetMaxOpenConns() int

	GetMaxIdleConns() int

	GetBlockBufferSize() uint8

	GetMaxCompressionBuffer() *pkl.DataSize

	GetEnableDebug() bool

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

	MigrationEngine string `pkl:"migrationEngine"`

	MigrationsPathClickhouse string `pkl:"migrationsPathClickhouse"`

	MaxExecutionTime *pkl.Duration `pkl:"maxExecutionTime"`

	CannMaxLifeTime *pkl.Duration `pkl:"cannMaxLifeTime"`

	DialTimeout *pkl.Duration `pkl:"dialTimeout"`

	MaxOpenConns int `pkl:"maxOpenConns"`

	MaxIdleConns int `pkl:"maxIdleConns"`

	BlockBufferSize uint8 `pkl:"blockBufferSize"`

	MaxCompressionBuffer *pkl.DataSize `pkl:"maxCompressionBuffer"`

	EnableDebug bool `pkl:"enableDebug"`

	TableName string `pkl:"tableName"`

	AsyncInsert string `pkl:"asyncInsert"`

	WaitForAsyncInsert string `pkl:"waitForAsyncInsert"`

	Name string `pkl:"name"`

	Inputs []string `pkl:"inputs"`

	Threads int `pkl:"threads"`
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

func (rcv *ClickhouseImpl) GetMigrationEngine() string {
	return rcv.MigrationEngine
}

func (rcv *ClickhouseImpl) GetMigrationsPathClickhouse() string {
	return rcv.MigrationsPathClickhouse
}

func (rcv *ClickhouseImpl) GetMaxExecutionTime() *pkl.Duration {
	return rcv.MaxExecutionTime
}

func (rcv *ClickhouseImpl) GetCannMaxLifeTime() *pkl.Duration {
	return rcv.CannMaxLifeTime
}

func (rcv *ClickhouseImpl) GetDialTimeout() *pkl.Duration {
	return rcv.DialTimeout
}

func (rcv *ClickhouseImpl) GetMaxOpenConns() int {
	return rcv.MaxOpenConns
}

func (rcv *ClickhouseImpl) GetMaxIdleConns() int {
	return rcv.MaxIdleConns
}

func (rcv *ClickhouseImpl) GetBlockBufferSize() uint8 {
	return rcv.BlockBufferSize
}

func (rcv *ClickhouseImpl) GetMaxCompressionBuffer() *pkl.DataSize {
	return rcv.MaxCompressionBuffer
}

func (rcv *ClickhouseImpl) GetEnableDebug() bool {
	return rcv.EnableDebug
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

func (rcv *ClickhouseImpl) GetThreads() int {
	return rcv.Threads
}
