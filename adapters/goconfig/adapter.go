// Package goconfig adapts github.com/majiddarvishan/goconfig to fluxlog.
//
// It is a separate Go module so applications using fluxlog's core package do
// not inherit goconfig or its transitive dependencies.
package goconfig

import (
	"fmt"
	"strings"

	"github.com/majiddarvishan/fluxlog"
	config "github.com/majiddarvishan/goconfig"
)

// New creates a fluxlog logger from a goconfig node and registers the level
// field for runtime replacement. The supplied node must use the legacy logger
// shape documented in this package's README.
func New(manager *config.Manager, node *config.Node) (*fluxlog.Logger, error) {
	if manager == nil {
		return nil, fmt.Errorf("fluxlog/goconfig: manager is nil")
	}

	loggerConfig, levelNode, err := Parse(node)
	if err != nil {
		return nil, err
	}

	logger, err := fluxlog.New(loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("fluxlog/goconfig: create logger: %w", err)
	}

	if err := BindLevel(manager, levelNode, logger); err != nil {
		_ = logger.Close()
		return nil, err
	}
	return logger, nil
}

// BindLevel connects a goconfig string node to a logger's runtime level. When
// a replacement is invalid, the callback returns an error and goconfig rolls
// the configuration change back.
func BindLevel(manager *config.Manager, levelNode *config.Node, logger *fluxlog.Logger) error {
	if manager == nil {
		return fmt.Errorf("fluxlog/goconfig: manager is nil")
	}
	if levelNode == nil {
		return fmt.Errorf("fluxlog/goconfig: level node is nil")
	}
	if logger == nil {
		return fmt.Errorf("fluxlog/goconfig: logger is nil")
	}

	if err := applyLevel(levelNode, logger); err != nil {
		return err
	}
	if err := manager.OnReplace(levelNode, func(replacement *config.Node) error {
		return applyLevel(replacement, logger)
	}); err != nil {
		return fmt.Errorf("fluxlog/goconfig: register level replacement: %w", err)
	}
	return nil
}

func applyLevel(node *config.Node, logger *fluxlog.Logger) error {
	level, err := node.GetString()
	if err != nil {
		return fmt.Errorf("fluxlog/goconfig: read level: %w", err)
	}
	if err := logger.SetLevelString(level); err != nil {
		return fmt.Errorf("fluxlog/goconfig: apply level: %w", err)
	}
	return nil
}

// Parse maps the original logger configuration shape to fluxlog.Config. It
// returns the level node separately so callers can choose whether to bind
// runtime updates with BindLevel.
func Parse(node *config.Node) (fluxlog.Config, *config.Node, error) {
	if node == nil {
		return fluxlog.Config{}, nil, fmt.Errorf("fluxlog/goconfig: config node is nil")
	}

	outputMode, err := requiredString(node, "output_mode")
	if err != nil {
		return fluxlog.Config{}, nil, err
	}
	levelNode, err := node.At("level")
	if err != nil {
		return fluxlog.Config{}, nil, fieldError("level", err)
	}
	levelText, err := levelNode.GetString()
	if err != nil {
		return fluxlog.Config{}, nil, fieldError("level", err)
	}
	level, err := fluxlog.ParseLevel(levelText)
	if err != nil {
		return fluxlog.Config{}, nil, fieldError("level", err)
	}

	result := fluxlog.Config{
		Level:     level,
		Timestamp: true,
		Caller:    true,
	}

	switch strings.ToLower(strings.TrimSpace(outputMode)) {
	case "console":
		result.Console = defaultConsole()
	case "file", "both":
		fileName, err := requiredString(node, "file_name")
		if err != nil {
			return fluxlog.Config{}, nil, err
		}
		maxFileSize, err := requiredInt(node, "max_file_size")
		if err != nil {
			return fluxlog.Config{}, nil, err
		}
		maxFiles, err := requiredInt(node, "max_files")
		if err != nil {
			return fluxlog.Config{}, nil, err
		}
		result.File = &fluxlog.FileConfig{
			Path:       fileName,
			Format:     fluxlog.JSONFormat,
			MaxSizeMB:  maxFileSize,
			MaxBackups: maxFiles,
			MaxAgeDays: 28,
		}
		if strings.EqualFold(strings.TrimSpace(outputMode), "both") {
			result.Console = defaultConsole()
		}
	default:
		return fluxlog.Config{}, nil, fieldError(
			"output_mode",
			fmt.Errorf("unsupported value %q; expected console, file, or both", outputMode),
		)
	}

	return result, levelNode, nil
}

func defaultConsole() *fluxlog.ConsoleConfig {
	return &fluxlog.ConsoleConfig{
		Format: fluxlog.ConsoleFormat,
		Color:  fluxlog.AutoColor,
	}
}

func requiredString(node *config.Node, name string) (string, error) {
	value, err := node.GetString(name)
	if err != nil {
		return "", fieldError(name, err)
	}
	return value, nil
}

func requiredInt(node *config.Node, name string) (int, error) {
	value, err := node.GetInt(name)
	if err != nil {
		return 0, fieldError(name, err)
	}
	return value, nil
}

func fieldError(name string, err error) error {
	return fmt.Errorf("fluxlog/goconfig: field %q: %w", name, err)
}
