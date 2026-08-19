package fluxlog

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

// Level is the minimum severity emitted by a Logger.
type Level string

const (
	TraceLevel    Level = "trace"
	DebugLevel    Level = "debug"
	InfoLevel     Level = "info"
	WarnLevel     Level = "warn"
	ErrorLevel    Level = "error"
	FatalLevel    Level = "fatal"
	PanicLevel    Level = "panic"
	DisabledLevel Level = "disabled"
)

// ParseLevel parses a level name without changing any process-global state.
func ParseLevel(value string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "trace":
		return TraceLevel, nil
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "panic":
		return PanicLevel, nil
	case "disabled", "off":
		return DisabledLevel, nil
	default:
		return "", fmt.Errorf("fluxlog: unsupported level %q", value)
	}
}

func (level Level) zerologLevel() (zerolog.Level, error) {
	parsed, err := ParseLevel(string(level))
	if err != nil {
		return zerolog.NoLevel, err
	}

	switch parsed {
	case TraceLevel:
		return zerolog.TraceLevel, nil
	case DebugLevel:
		return zerolog.DebugLevel, nil
	case InfoLevel:
		return zerolog.InfoLevel, nil
	case WarnLevel:
		return zerolog.WarnLevel, nil
	case ErrorLevel:
		return zerolog.ErrorLevel, nil
	case FatalLevel:
		return zerolog.FatalLevel, nil
	case PanicLevel:
		return zerolog.PanicLevel, nil
	case DisabledLevel:
		return zerolog.Disabled, nil
	default:
		return zerolog.NoLevel, fmt.Errorf("fluxlog: unsupported level %q", level)
	}
}
