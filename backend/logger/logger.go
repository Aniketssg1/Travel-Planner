package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Log is the global application logger.
var Log zerolog.Logger

// Init initialises the global zerolog logger.
// In "debug" mode it uses a pretty console writer; otherwise JSON.
func Init(ginMode string) {
	level := zerolog.InfoLevel
	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		if parsed, err := zerolog.ParseLevel(lvl); err == nil {
			level = parsed
		}
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339

	if ginMode == "debug" {
		Log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}).
			With().Timestamp().Caller().Logger()
	} else {
		Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}
