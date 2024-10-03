package pipelaner

import (
	"io"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type loggerKey int

const (
	logKey loggerKey = iota
)

func initLogger(cfg *Config) (*zerolog.Logger, error) {
	var writers []io.Writer

	if cfg.LogLevel == "" {
		cfg.LogLevel = zerolog.LevelDebugValue
	}
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	if cfg.EnableConsole {
		if cfg.LogFormat == LogFormatJSON {
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339Nano,
			})
		} else {
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: time.RFC3339Nano,
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

func newRollingFile(cfg *Config) io.Writer {
	return &lumberjack.Logger{
		Filename:   path.Join(cfg.FileDirectory, cfg.FileName),
		MaxBackups: cfg.FileMaxBackups,
		MaxSize:    cfg.FileMaxSize,
		MaxAge:     cfg.FileMaxAge,
		Compress:   cfg.FileCompress,
		LocalTime:  cfg.FileLocalTime,
	}
}

func NewLogger() zerolog.Logger {
	lg := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).Level(zerolog.InfoLevel).
		With().
		Timestamp()
	l := lg.Logger()
	return l
}
