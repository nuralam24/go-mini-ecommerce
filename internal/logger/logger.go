package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Log zerolog.Logger

func Init(environment string) {
	zerolog.TimeFieldFormat = time.RFC3339

	if environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	Log = log.Logger.With().
		Str("service", "go-ecommerce").
		Logger()
}

func Get() *zerolog.Logger {
	return &Log
}
