package migrator

import (
	"fmt"

	config "github.com/pipelane/pipelaner/gen/pipelaner"
	"github.com/pipelane/pipelaner/gen/settings/migrations"
	"github.com/pipelane/pipelaner/gen/settings/migrations/driver"
	"github.com/pipelane/pipelaner/internal/logger"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type Migrator struct {
	logger     *zerolog.Logger
	cfg        migrations.Config
	migrations map[driver.Driver]MigrationInterface
}

func NewMigrator(
	cfg *config.Pipelaner,
) (*Migrator, error) {
	l, err := logger.NewLoggerWithCfg(cfg.Settings.Logger)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	c := cfg.Settings.Migrations
	migrators := map[driver.Driver]MigrationInterface{}
	for _, m := range *c.Migrations {
		switch m.GetDriver() {
		case driver.Clickhouse:
			c, ok := m.(migrations.Clickhouse)
			if !ok {
				return nil, fmt.Errorf("invalid migration type: %T", m)
			}
			migrators[m.GetDriver()] = NewMigrateClick(c)
		case driver.Empty:
			return nil, fmt.Errorf("empty migrations not supported")
		}
	}
	return &Migrator{
		cfg:        *c,
		migrations: migrators,
		logger:     l,
	}, nil
}
func (m *Migrator) Migrate() error {
	gr := errgroup.Group{}
	for _, mgr := range *m.cfg.Migrations {
		m.logger.Info().Msgf("Starting migration: %s", mgr.GetDriver())
		migrator := m.migrations[mgr.GetDriver()]
		gr.Go(func() error {
			err := migrator.Run(mgr.GetPath())
			if err != nil {
				m.logger.Error().Err(err).Msgf("Failed migration: %s", mgr.GetDriver())
				return err
			}
			return nil
		})
	}
	return gr.Wait()
}
