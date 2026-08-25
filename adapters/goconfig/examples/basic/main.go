package main

import (
	"log"

	fluxgoconfig "github.com/majiddarvishan/fluxlog/adapters/goconfig"
	config "github.com/majiddarvishan/goconfig"
)

const configJSON = `{
  "logger": {
    "output_mode": "both",
    "level": "info",
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
      "required": ["output_mode", "level", "file_name", "max_file_size", "max_files"]
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

	logger.Info().Str("component", "example").Msg("configured through goconfig")
}
