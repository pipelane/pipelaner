// Code generated from Pkl module `com.pipelaner.settings.migrations.config`. DO NOT EDIT.
package migrations

import "github.com/pipelane/pipelaner/gen/settings/migrations/driver"

type Migration interface {
	GetDriver() driver.Driver

	GetPath() string
}
