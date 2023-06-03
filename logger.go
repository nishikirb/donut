package donut

import (
	"io"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	logger     zerolog.Logger
	initLogger sync.Once
)

func InitLogger(out io.Writer, isDebug bool) error {
	initLogger.Do(func() {
		writer := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.Out = out
			w.TimeFormat = time.RFC3339
		})
		if isDebug {
			logger = log.Output(writer).Level(zerolog.InfoLevel)
		} else {
			logger = log.Output(writer).Level(zerolog.Disabled)
		}
	})
	return nil
}

func GetLogger() zerolog.Logger {
	return logger
}
