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

func NewLoggerWithCfg(cfg *logger.LoggerConfig) (*zerolog.Logger, error) {
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
				Out: os.Stderr,
			})
		}
	}
	if cfg.EnableFile {
		writers = append(writers, newRollingFile(cfg))
	}
	mw := io.MultiWriter(writers...)
	logger := zerolog.
		New(mw).
		Level(level).
		With().Timestamp().
		Logger()
	return &logger, nil
}

func newRollingFile(cfg *logger.LoggerConfig) io.Writer {
	return &lumberjack.Logger{
		Filename:   path.Join(*cfg.FileDirectory, *cfg.FileName),
		MaxBackups: *cfg.FileMaxBackups,
		MaxSize:    int(cfg.FileMaxSize.ToUnit(pkl.Megabytes).Value),
		MaxAge:     *cfg.FileMaxAge,
		Compress:   *cfg.FileCompress,
		LocalTime:  *cfg.FileLocalFormat,
	}
}
