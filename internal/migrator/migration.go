// Author: Alexey Khokhlov
//

package migrator

import (
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pipelane/pipelaner/gen/settings/migrations"
)

type MigrationInterface interface {
	Run(migrationsDir string) error
}

type ClickhouseMigration interface {
	MigrationInterface
}

type Click struct {
	cfg migrations.Clickhouse
}

func NewMigrateClick(cfg migrations.Clickhouse) *Click {
	return &Click{cfg: cfg}
}

func (m *Click) Run(migrationsDir string) error {
	p := &ClickHouse{}
	a := strings.Split(m.cfg.GetCredentials().Address, ":")
	addr := clickhouseConnectionString(
		m.cfg.GetCredentials().User,
		m.cfg.GetCredentials().Password,
		a[0],
		a[1],
		m.cfg.GetCredentials().Database,
		m.cfg.GetEngine(),
	)
	d, err := p.Open(addr)
	if err != nil {
		return err
	}
	migration, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		m.cfg.GetCredentials().Database, d)
	if err != nil {
		return err
	}
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	if err := d.Close(); err != nil {
		return err
	}
	return nil
}

func clickhouseConnectionString(user, password, host, port, db string, engine *string) string {
	if engine != nil {
		return fmt.Sprintf(
			"clickhouse://%s:%s@%v:%v/%s?x-multi-statement=true&x-migrations-table-engine=%v&debug=false",
			user, password, host, port, db, *engine)
	}
	return fmt.Sprintf(
		"clickhouse://%s:%s@%v:%v/%s?x-multi-statement=true&debug=false",
		user, password, host, port, db)
}
