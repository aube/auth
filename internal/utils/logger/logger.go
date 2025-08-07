// Package logger provides centralized logging configuration.
package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init configures global zerolog settings.
// logLevel: "debug" or empty (info)
// Configures:
//   - Console output with timestamps
//   - Debug or info log level
//   - Structured logging
func Init(logLevel string) {
	level := zerolog.InfoLevel

	if logLevel == "debug" {
		level = zerolog.DebugLevel
	}

	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05"}
	zerolog.SetGlobalLevel(level)
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}

// Get returns the configured global logger instance.
// Returns: zerolog.Logger
// Usage:
//
//	logger.Get().Info().Msg("message")
//	logger.Get().Debug().Str("key", "value").Msg("debug")
func Get() zerolog.Logger {
	return log.Logger
}
