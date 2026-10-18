package main

import (
	"log"

	"github.com/majiddarvishan/fluxlog"
	fluxgoconfig "github.com/majiddarvishan/fluxlog/adapters/goconfig"
	config "github.com/majiddarvishan/goconfig"
)

const configJSON = `{
  "logger": {
    "output_mode": "both",
    "level": "info",
    "service": "gateway",
    "caller_max_length": 15,
    "file_name": "./logs/application.log",
    "max_file_size": 100,
    "max_files": 10
  }
}`

const schemaJSON = `{
  "type": "object",
  "required": ["logger"],
  "properties": {
    "logger": {
      "type": "object",
      "required": ["output_mode", "level", "file_name", "max_file_size", "max_files"],
      "properties": {
        "output_mode": {"type": "string"},
        "level": {"type": "string"},
        "service": {"type": "string"},
        "caller_max_length": {"type": "integer"},
        "file_name": {"type": "string"},
        "max_file_size": {"type": "integer"},
        "max_files": {"type": "integer"}
      }
    }
  }
}`

func main() {
	source, err := config.NewStrSource(configJSON, schemaJSON)
	if err != nil {
		log.Fatal(err)
	}
	manager, err := config.NewManager(source)
	if err != nil {
		log.Fatal(err)
	}
	loggerNode, err := manager.Config().At("logger")
	if err != nil {
		log.Fatal(err)
	}

	// New parses this node, creates a Fluxlog instance and binds its level
	// to successfully committed goconfig replacements.
	logger, err := fluxgoconfig.New(manager, loggerNode)
	if err != nil {
		log.Fatal(err)
	}
	previous := fluxlog.SetDefault(logger)
	defer func() {
		fluxlog.SetDefault(previous) // Restore before closing the owned logger.
		_ = logger.Close()
	}()

	fluxlog.Info("Web server starting on %s", "0.0.0.0:4567")
	fluxlog.Debug("Hidden while the configured level is info")

	// Replace goes through goconfig's validation and persistence callbacks.
	// There is no need to call logger.SetLevelString here.
	if err := manager.Replace("/logger/level", "debug"); err != nil {
		log.Printf("change configured log level: %v", err)
		return
	}
	fluxlog.Debug("Debug logging enabled by goconfig")
	fluxlog.Info("Processed %d requests", 42)
	logger.Info().Str("user_id", "123").Msg("User authenticated")

	// A rejected replacement must leave the active level unchanged.
	if err := manager.Replace("/logger/level", "invalid-level"); err != nil {
		fluxlog.Warn("Invalid level rejected: %v", err)
	} else {
		log.Print("unexpected success replacing log level with an invalid value")
		return
	}
	fluxlog.Debug("Debug is still enabled after the rejected replacement")
}
