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

	logger, err := fluxgoconfig.New(manager, loggerNode)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	previous := fluxlog.SetDefault(logger)
	defer fluxlog.SetDefault(previous)

	fluxlog.Info("Web server starting on %s", "0.0.0.0:4567")
}
