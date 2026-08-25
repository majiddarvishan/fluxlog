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
		Console: &fluxlog.ConsoleConfig{
			Format: fluxlog.ConsoleFormat,
			Color:  fluxlog.AutoColor,
		},
		File: &fluxlog.FileConfig{
			Path:       "./logs/application.log",
			Format:     fluxlog.JSONFormat,
			MaxSizeMB:  100,
			MaxBackups: 10,
			MaxAgeDays: 28,
			Compress:   true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	logger.Info().Int("port", 8080).Msg("service started")
}
