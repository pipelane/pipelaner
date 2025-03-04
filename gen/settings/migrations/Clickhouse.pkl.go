// Code generated from Pkl module `com.pipelaner.settings.migrations.config`. DO NOT EDIT.
package migrations

import (
	"github.com/pipelane/pipelaner/gen/settings/migrations/driver"
	"github.com/pipelane/pipelaner/gen/source/common"
)

type Clickhouse interface {
	Migration

	GetCredentials() *common.ChCredentials

	GetEngine() string
}

var _ Clickhouse = (*ClickhouseImpl)(nil)

type ClickhouseImpl struct {
	Driver driver.Driver `pkl:"driver"`

	Path string `pkl:"path"`

	Credentials *common.ChCredentials `pkl:"credentials"`

	Engine string `pkl:"engine"`
}

func (rcv *ClickhouseImpl) GetDriver() driver.Driver {
	return rcv.Driver
}

func (rcv *ClickhouseImpl) GetPath() string {
	return rcv.Path
}

func (rcv *ClickhouseImpl) GetCredentials() *common.ChCredentials {
	return rcv.Credentials
}

func (rcv *ClickhouseImpl) GetEngine() string {
	return rcv.Engine
}
