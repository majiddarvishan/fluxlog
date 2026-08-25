package fluxlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/rs/zerolog"
)

func newConsoleWriter(
	output io.Writer,
	color ColorMode,
	timeFormat string,
	callerMaxLength int,
) zerolog.ConsoleWriter {
	return zerolog.ConsoleWriter{
		Out:           output,
		NoColor:       !colorEnabled(color, output),
		TimeFormat:    timeFormat,
		PartsOrder:    []string{zerolog.TimestampFieldName, zerolog.LevelFieldName, zerolog.CallerFieldName, zerolog.MessageFieldName},
		FieldsExclude: []string{"service"},
		FormatCaller:  callerFormatter(callerMaxLength),
		FormatPrepare: prepareConsoleEvent,
	}
}

func prepareConsoleEvent(event map[string]any) error {
	service, ok := event["service"].(string)
	if !ok || strings.TrimSpace(service) == "" {
		return nil
	}

	prefix := "(" + capitalizeFirst(service) + ")"
	message, _ := event[zerolog.MessageFieldName].(string)
	if message == "" {
		event[zerolog.MessageFieldName] = prefix
	} else {
		event[zerolog.MessageFieldName] = prefix + " " + message
	}
	return nil
}

func callerFormatter(maxLength int) zerolog.Formatter {
	return func(value any) string {
		caller, ok := value.(string)
		if !ok || caller == "" {
			return ""
		}
		if workingDirectory, err := os.Getwd(); err == nil {
			if relative, err := filepath.Rel(workingDirectory, caller); err == nil {
				caller = relative
			}
		}
		return "[ " + truncateLeft(caller, maxLength) + " ]"
	}
}

func truncateLeft(value string, maxLength int) string {
	runes := []rune(value)
	if len(runes) <= maxLength {
		return value
	}
	if maxLength == 1 {
		return "…"
	}
	return "…" + string(runes[len(runes)-maxLength+1:])
}

func capitalizeFirst(value string) string {
	value = strings.TrimSpace(value)
	first, size := utf8.DecodeRuneInString(value)
	if first == utf8.RuneError && size == 0 {
		return value
	}
	return fmt.Sprintf("%c%s", unicode.ToUpper(first), value[size:])
}
