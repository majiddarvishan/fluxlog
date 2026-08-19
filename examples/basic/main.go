package main

import (
	"log"

	"github.com/majiddarvishan/fluxlog"
)

func main() {
	logger, err := fluxlog.New(fluxlog.Config{
		Level:     fluxlog.InfoLevel,
		Service:   "example-service",
		Timestamp: true,
		Caller:    true,
		Console: &fluxlog.ConsoleConfig{
			Format: fluxlog.ConsoleFormat,
			Color:  fluxlog.AutoColor,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	logger.Info().Str("request_id", "req-1").Msg("request received")

	if err := logger.SetLevelString("debug"); err != nil {
		log.Fatal(err)
	}
	logger.Debug().Msg("debug logging enabled at runtime")
}
