/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package clickhouse

import (
	"github.com/pipelane/pipelaner"
	"time"
)

type ClickhouseConfig struct {
	Address                  string        `pipelane:"address"`
	User                     string        `pipelane:"user"`
	Password                 string        `pipelane:"password"`
	Database                 string        `pipelane:"database"`
	MigrationEngine          string        `pipelane:"migration_engine"`
	MigrationsPathClickhouse string        `pipelane:"migrations_path_clickhouse"`
	MaxExecutionTime         time.Duration `pipelane:"max_execution_time"`
	ConnMaxLifetime          time.Duration `pipelane:"conn_max_lifetime"`
	DialTimeout              time.Duration `pipelane:"dial_timeout"`
	MaxOpenConns             int           `pipelane:"max_open_conns"`
	MaxIdleConns             int           `pipelane:"max_idle_conns"`
	BlockBufferSize          uint8         `pipelane:"block_buffer_size"`
	MaxCompressionBuffer     string        `pipelane:"max_compression_buffer"`
	EnableDebug              bool          `pipelane:"enable_debug"`
	TableName                string        `pipelane:"table_name"`
	pipelaner.Internal
}
