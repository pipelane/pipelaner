// Author: Alexey Khokhlov
//

package migrator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"

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
	e := m.cfg.GetEngine()
	addr := clickhouseConnectionString(
		m.cfg.GetCredentials().User,
		m.cfg.GetCredentials().Password,
		a[0],
		a[1],
		m.cfg.GetCredentials().Database,
		&e,
		m.cfg.GetClusterName(),
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
	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	if err = d.Close(); err != nil {
		return err
	}
	return nil
}

func clickhouseConnectionString(user, password, host, port, db string, engine, cluster *string) string {
	url := fmt.Sprintf(
		"clickhouse://%s:%s@%v:%v/%s?x-multi-statement=true&debug=false",
		user, password, host, port, db)
	if engine != nil {
		url += fmt.Sprintf("&x-migrations-table-engine=%s", *engine)
	}
	if cluster != nil {
		url += fmt.Sprintf("&x-cluster-name=%s", *cluster)
	}
	return url
}
