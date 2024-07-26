package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger interface {
	Trace() *zerolog.Event
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	WithLevel(level zerolog.Level) *zerolog.Event

	With() zerolog.Context
}

func NewLogger() Logger {
	var output io.Writer = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	if os.Getenv("GO_ENV") != "development" {
		output = os.Stderr
	}

	logger := zerolog.New(output).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	return &logger
}
