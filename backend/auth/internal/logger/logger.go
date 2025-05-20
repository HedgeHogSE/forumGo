package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func GetLogger(component string) zerolog.Logger {
	return log.With().Str("component", component).Logger()
}
