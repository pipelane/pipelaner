// Code generated from Pkl module `pipelaner.source.sinks`. DO NOT EDIT.
package sink

type Clickhouse interface {
	Sink

	GetAddress() string

	GetPassword() string

	GetDatabase() string

	GetMigrationEngine() string

	GetMigrationsPathClickhouse() string

	GetMaxExecutionTime() string
}

var _ Clickhouse = (*ClickhouseImpl)(nil)

type ClickhouseImpl struct {
	SourceName string `pkl:"sourceName"`

	Address string `pkl:"address"`

	Password string `pkl:"password"`

	Database string `pkl:"database"`

	MigrationEngine string `pkl:"migrationEngine"`

	MigrationsPathClickhouse string `pkl:"migrationsPathClickhouse"`

	MaxExecutionTime string `pkl:"maxExecutionTime"`

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

func (rcv *ClickhouseImpl) GetMaxExecutionTime() string {
	return rcv.MaxExecutionTime
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
