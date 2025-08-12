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

func NewLoggerWithCfg(cfg logger.Config) (*zerolog.Logger, error) {
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
			if cfg.FileParams != nil {
				writers = append(writers, newRollingFile(cfg))
			}
		} else {
			writer := zerolog.ConsoleWriter{
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
			}
			writer.Out = os.Stdout
			writers = append(writers, writer)
			if cfg.FileParams != nil {
				writer.Out = newRollingFile(cfg)
				writers = append(writers, writer)
			}
		}
	}
	mw := io.MultiWriter(writers...)
	l := zerolog.
		New(mw).
		Level(level).
		With().Timestamp().
		Logger()
	return &l, nil
}

func newRollingFile(cfg logger.Config) io.Writer {
	maxSize := cfg.FileParams.MaxSize.ToUnit(pkl.Megabytes)
	return &lumberjack.Logger{
		Filename:   path.Join(cfg.FileParams.Directory, cfg.FileParams.Name),
		MaxBackups: cfg.FileParams.MaxBackups,
		MaxSize:    int(maxSize.Value),
		MaxAge:     cfg.FileParams.MaxAge,
		Compress:   cfg.FileParams.Compress,
		LocalTime:  cfg.FileParams.LocalFormat,
	}
}
