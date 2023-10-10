package pipelane

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

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
