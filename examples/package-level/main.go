package main

import (
	"log"

	"github.com/majiddarvishan/fluxlog"
)

func main() {
	logger, err := fluxlog.New(fluxlog.Config{
		Level:     fluxlog.DebugLevel,
		Service:   "gateway",
		Timestamp: true,
		Caller:    true,
		Console: &fluxlog.ConsoleConfig{
			Format:          fluxlog.ConsoleFormat,
			Color:           fluxlog.AutoColor,
			CallerMaxLength: 15,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	previous := fluxlog.SetDefault(logger)
	defer fluxlog.SetDefault(previous)

	fluxlog.Info("server started")
	fluxlog.Debug("listening on port %d", 8080)
	fluxlog.Warn("retry %d of %d", 1, 3)
}
