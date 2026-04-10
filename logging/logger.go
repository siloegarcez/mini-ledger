package logging

import (
	"mini-ledger/config"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Configure(conf config.Config) {
	if conf.LogLevel < 0 || conf.LogLevel > 7 {
		log.Warn().Int("level", conf.LogLevel).Msg("Invalid log level, defaulting to debug")
		conf.LogLevel = 0
	}
	baseLogger := log.With().
		Timestamp().
		Caller().
		Stack().
		Str("env", conf.Environment).
		Str("version", conf.Version).
		Logger().Level(zerolog.Level(conf.LogLevel)) //nolint:gosec

	if conf.Environment != config.EnvDev {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = baseLogger.Output(os.Stderr)
		return
	}

	log.Logger = baseLogger.Output(zerolog.ConsoleWriter{Out: os.Stderr}) // nolint:exhaustruct
}
