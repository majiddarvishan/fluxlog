package main

import (
	"log"

	"github.com/majiddarvishan/fluxlog"
)

func main() {
	logger, err := fluxlog.New(fluxlog.Config{
		Level:     fluxlog.InfoLevel,
		Service:   "gateway",
		Timestamp: true,
		Caller:    true,
		Console: &fluxlog.ConsoleConfig{
			Format:          fluxlog.ConsoleFormat,
			Color:           fluxlog.AutoColor,
			CallerMaxLength: 15,
			TimeFormat:      "2006-01-02 15:04:05",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Package-level Info/Debug/Warn calls now use this configured instance.
	// SetDefault returns the former default; it does not close either logger.
	previous := fluxlog.SetDefault(logger)
	defer func() {
		fluxlog.SetDefault(previous) // Stop using this logger as the default.
		_ = logger.Close()           // Then release its owned outputs.
	}()

	fluxlog.Info("Application started")
	fluxlog.Info("Listening on port %d", 8080)

	// Debug is suppressed until this instance's level is changed.
	fluxlog.Debug("This debug message is hidden")
	if err := logger.SetLevelString("debug"); err != nil {
		log.Printf("change log level: %v", err)
		return
	}
	fluxlog.Debug("Debug logging enabled")
	fluxlog.Warn("Retry %d of %d", 1, 3)

	// Use the instance API when a log event needs structured fields.
	logger.Info().Str("user_id", "123").Msg("User authenticated")
}
