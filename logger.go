package pipelaner

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
)

type loggerKey int

const (
	logKey loggerKey = iota
)

func initLogger(cfg *config) (*zerolog.Logger, error) {
	var writers []io.Writer

	if cfg.LogLevel == "" {
		cfg.LogLevel = zerolog.LevelDebugValue
	}
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	if cfg.EnableConsole {
		writers = append(writers, zerolog.ConsoleWriter{
			Out: os.Stderr,
		})
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

func newRollingFile(cfg *config) io.Writer {
	return &lumberjack.Logger{
		Filename:   path.Join(cfg.FileDirectory, cfg.FileName),
		MaxBackups: cfg.FileMaxBackups,
		MaxSize:    cfg.FileMaxSize,
		MaxAge:     cfg.FileMaxAge,
		Compress:   cfg.FileCompress,
		LocalTime:  cfg.FileLocalTime,
	}
}
