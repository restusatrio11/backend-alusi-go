package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Log zerolog.Logger

// InitLogger initializes the global Zerolog logger
func InitLogger(isDebug bool, appEnv string) {
	zerolog.TimeFieldFormat = time.RFC3339

	if appEnv == "development" || isDebug {
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}
		Log = zerolog.New(output).With().Timestamp().Caller().Logger()
	} else {
		Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	if isDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = Log
}
