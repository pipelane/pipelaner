/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/apple/pkl-go/pkl"
	"github.com/pipelane/pipelaner/gen/settings/logger"
	"github.com/pipelane/pipelaner/gen/settings/logger/logformat"
	"github.com/pipelane/pipelaner/gen/settings/logger/loglevel"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLoggerWithCfg(cfg *logger.Config) (*zerolog.Logger, error) {
	var writers []io.Writer

	if cfg.LogLevel == "" {
		cfg.LogLevel = loglevel.Info
	}
	level, err := zerolog.ParseLevel(cfg.LogLevel.String())
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano
	if cfg.EnableConsole {
		if cfg.LogFormat == logformat.Json {
			writers = append(writers, os.Stdout)
		} else {
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339Nano,
				FieldsOrder: []string{
					"time",
					zerolog.LevelFieldName,
					zerolog.ErrorFieldName,
					"type",
					"lane_name",
					"source",
					zerolog.MessageFieldName,
					zerolog.CallerFieldName,
				},
			})
		}
	}
	if cfg.FileParams != nil {
		writers = append(writers, newRollingFile(cfg))
	}
	mw := io.MultiWriter(writers...)
	l := zerolog.
		New(mw).
		Level(level).
		With().Timestamp().
		Logger()
	return &l, nil
}

func newRollingFile(cfg *logger.Config) io.Writer {
	val := cfg.FileParams.MaxSize
	orig := (val.Value * float64(val.Unit)) / pkl.Megabytes

	return &lumberjack.Logger{
		Filename:   path.Join(cfg.FileParams.Directory, cfg.FileParams.Name),
		MaxBackups: cfg.FileParams.MaxBackups,
		MaxSize:    int(orig),
		MaxAge:     cfg.FileParams.MaxAge,
		Compress:   cfg.FileParams.Compress,
		LocalTime:  cfg.FileParams.LocalFormat,
	}
}
