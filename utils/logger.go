package utils

import (
	"os"

	"github.com/rs/zerolog"
)

var BaseLogger zerolog.Logger

func init() {
	BaseLogger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
}

func ChildLogger(component string) zerolog.Logger {
	return BaseLogger.With().Str("component", component).Logger()
}
