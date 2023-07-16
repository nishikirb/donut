package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	logger zerolog.Logger = New(os.Stdout, false)
	once   sync.Once
)

// New creates a new logger.
func New(out io.Writer, verbose bool) zerolog.Logger {
	var l zerolog.Logger
	writer := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = out
		w.TimeFormat = time.RFC3339
	})
	if verbose {
		l = log.Output(writer).Level(zerolog.InfoLevel)
	} else {
		l = log.Output(writer).Level(zerolog.Disabled)
	}

	return l
}

// Init initializes the global logger.
func Init(out io.Writer, verbose bool) {
	once.Do(func() {
		logger = New(out, verbose)
	})
}

// Info returns a new event with the info level.
func Info() *zerolog.Event {
	return logger.Info()
}
