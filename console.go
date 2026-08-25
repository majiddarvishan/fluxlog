package fluxlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/rs/zerolog"
)

const (
	consoleServiceField = "service"
	ansiCyan            = "\x1b[36m"
	ansiReset           = "\x1b[0m"
)

func newConsoleWriter(
	output io.Writer,
	color ColorMode,
	timeFormat string,
	callerMaxLength int,
) zerolog.ConsoleWriter {
	writer := zerolog.ConsoleWriter{Out: output}
	noColor := !colorEnabled(color, output)
	writer.NoColor = noColor
	writer.TimeFormat = timeFormat
	writer.PartsOrder = []string{zerolog.TimestampFieldName, zerolog.LevelFieldName, zerolog.CallerFieldName, consoleServiceField, zerolog.MessageFieldName}
	writer.FieldsExclude = []string{consoleServiceField}
	writer.FormatTimestamp = timestampFormatter(timeFormat)
	writer.FormatCaller = callerFormatter(callerMaxLength)
	writer.FormatPartValueByName = consolePartFormatter(noColor)
	return writer
}

func timestampFormatter(timeFormat string) zerolog.Formatter {
	return func(value any) string {
		timestamp, ok := value.(string)
		if !ok || timestamp == "" {
			return ""
		}
		parsed, err := time.Parse(zerolog.TimeFieldFormat, timestamp)
		if err != nil {
			return timestamp
		}
		return parsed.Format(timeFormat)
	}
}

func consolePartFormatter(noColor bool) zerolog.FormatterByFieldName {
	return func(value any, field string) string {
		if field != consoleServiceField {
			return fmt.Sprint(value)
		}
		service, ok := value.(string)
		if !ok || strings.TrimSpace(service) == "" {
			return ""
		}
		formatted := "(" + capitalizeFirst(service) + ")"
		if noColor || os.Getenv("NO_COLOR") != "" {
			return formatted
		}
		return ansiCyan + formatted + ansiReset
	}
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
