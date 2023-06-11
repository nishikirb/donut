package donut

import (
	"io"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger(out io.Writer, verbose bool) zerolog.Logger {
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
