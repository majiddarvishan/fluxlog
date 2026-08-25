package fluxlog

import (
	"fmt"
	"sync/atomic"

	"github.com/rs/zerolog"
)

var packageLogger atomic.Pointer[Logger]

func init() {
	logger, err := New(DefaultConfig())
	if err != nil {
		panic(err)
	}
	packageLogger.Store(logger)
}

// Default returns the Logger used by the package-level logging functions.
func Default() *Logger {
	return packageLogger.Load()
}

// SetDefault replaces the Logger used by the package-level logging functions
// and returns the previous Logger. Logger lifecycle remains owned by the
// caller; SetDefault does not close either Logger.
//
// SetDefault panics when logger is nil.
func SetDefault(logger *Logger) *Logger {
	if logger == nil {
		panic("fluxlog: nil default logger")
	}
	return packageLogger.Swap(logger)
}

// Trace writes a trace message through the default Logger.
func Trace(values ...any) {
	writeMessage(Default().Trace(), values)
}

// Debug writes a debug message through the default Logger.
func Debug(values ...any) {
	writeMessage(Default().Debug(), values)
}

// Info writes an info message through the default Logger.
func Info(values ...any) {
	writeMessage(Default().Info(), values)
}

// Warn writes a warning message through the default Logger.
func Warn(values ...any) {
	writeMessage(Default().Warn(), values)
}

// Error writes an error message through the default Logger.
func Error(values ...any) {
	writeMessage(Default().Error(), values)
}

// Fatal writes a fatal message and terminates the process with exit status 1
// when fatal logging is enabled.
func Fatal(values ...any) {
	writeMessage(Default().Fatal(), values)
}

// Panic writes a panic message and then panics when panic logging is enabled.
func Panic(values ...any) {
	writeMessage(Default().Panic(), values)
}

func writeMessage(event *zerolog.Event, values []any) {
	if event == nil {
		return
	}
	if len(values) == 0 {
		event.Msg("")
		return
	}
	if len(values) == 1 {
		if message, ok := values[0].(string); ok {
			event.Msg(message)
			return
		}
		event.Msg(fmt.Sprint(values[0]))
		return
	}
	if format, ok := values[0].(string); ok {
		event.Msgf(format, values[1:]...)
		return
	}
	event.Msg(fmt.Sprint(values...))
}
