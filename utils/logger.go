package utils

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var BaseLogger zerolog.Logger

func init() {
	logFile, err := os.OpenFile(filepath.Join(getProjectRoot(), "logs.json"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Fallback to stderr if file can't be opened
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)

	BaseLogger = zerolog.New(multi).With().Timestamp().Logger()
}

func getProjectRoot() string {
	dir, _ := os.Getwd()
	return dir
}

func ChildLogger(component string) zerolog.Logger {
	return BaseLogger.With().Str("component", component).Logger()
}
