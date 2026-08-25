package fluxlog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// OutputFormat controls how an output encodes log events.
type OutputFormat string

const (
	JSONFormat    OutputFormat = "json"
	ConsoleFormat OutputFormat = "console"
)

// ColorMode controls ANSI colors in console-formatted output.
type ColorMode string

const (
	AutoColor   ColorMode = "auto"
	AlwaysColor ColorMode = "always"
	NeverColor  ColorMode = "never"
)

// ConsoleConfig configures an optional stream output. The presence of this
// value enables the output. Writer defaults to os.Stdout.
type ConsoleConfig struct {
	Writer          io.Writer
	Format          OutputFormat
	Color           ColorMode
	TimeFormat      string
	CallerMaxLength int
}

// DefaultCallerMaxLength is the maximum number of visible characters used for
// a caller path in console-formatted output.
const DefaultCallerMaxLength = 15

// FileConfig configures an optional size-rotated file output. The presence of
// this value enables the output. Sizes are expressed in megabytes, matching
// lumberjack's API.
type FileConfig struct {
	Path       string
	Format     OutputFormat
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
	LocalTime  bool
}

// Config contains immutable construction settings. Runtime-safe changes are
// exposed as explicit methods on Logger instead of being tied to a config
// provider.
type Config struct {
	Level     Level
	Service   string
	Timestamp bool
	Caller    bool
	Console   *ConsoleConfig
	File      *FileConfig
}

// DefaultConfig returns a development-friendly console configuration.
func DefaultConfig() Config {
	return Config{
		Level:     InfoLevel,
		Timestamp: true,
		Console: &ConsoleConfig{
			Writer:          os.Stdout,
			Format:          ConsoleFormat,
			Color:           AutoColor,
			TimeFormat:      time.RFC3339,
			CallerMaxLength: DefaultCallerMaxLength,
		},
	}
}

func (config Config) normalized() (Config, error) {
	if strings.TrimSpace(string(config.Level)) == "" {
		config.Level = InfoLevel
	}
	if _, err := config.Level.zerologLevel(); err != nil {
		return Config{}, err
	}

	if config.Console == nil && config.File == nil {
		return Config{}, errors.New("fluxlog: at least one output is required")
	}

	if config.Console != nil {
		console := *config.Console
		if console.Writer == nil {
			console.Writer = os.Stdout
		}
		if console.Format == "" {
			console.Format = ConsoleFormat
		}
		if err := validateFormat(console.Format); err != nil {
			return Config{}, fmt.Errorf("fluxlog: console output: %w", err)
		}
		if console.Color == "" {
			console.Color = AutoColor
		}
		if err := validateColor(console.Color); err != nil {
			return Config{}, fmt.Errorf("fluxlog: console output: %w", err)
		}
		if console.TimeFormat == "" {
			console.TimeFormat = time.RFC3339
		}
		if console.CallerMaxLength < 0 {
			return Config{}, errors.New("fluxlog: console caller max length cannot be negative")
		}
		if console.CallerMaxLength == 0 {
			console.CallerMaxLength = DefaultCallerMaxLength
		}
		config.Console = &console
	}

	if config.File != nil {
		file := *config.File
		file.Path = strings.TrimSpace(file.Path)
		if file.Path == "" {
			return Config{}, errors.New("fluxlog: file output path is required")
		}
		if file.Format == "" {
			file.Format = JSONFormat
		}
		if err := validateFormat(file.Format); err != nil {
			return Config{}, fmt.Errorf("fluxlog: file output: %w", err)
		}
		if file.MaxSizeMB < 0 || file.MaxBackups < 0 || file.MaxAgeDays < 0 {
			return Config{}, errors.New("fluxlog: file rotation values cannot be negative")
		}
		if file.MaxSizeMB == 0 {
			file.MaxSizeMB = 100
		}
		config.File = &file
	}

	config.Service = strings.TrimSpace(config.Service)
	return config, nil
}

func validateFormat(format OutputFormat) error {
	switch format {
	case JSONFormat, ConsoleFormat:
		return nil
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

func validateColor(mode ColorMode) error {
	switch mode {
	case AutoColor, AlwaysColor, NeverColor:
		return nil
	default:
		return fmt.Errorf("unsupported color mode %q", mode)
	}
}
