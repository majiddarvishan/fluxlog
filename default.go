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
	writeMessage(Default().packageEvent(zerolog.TraceLevel), values)
}

// Debug writes a debug message through the default Logger.
func Debug(values ...any) {
	writeMessage(Default().packageEvent(zerolog.DebugLevel), values)
}

// Info writes an info message through the default Logger.
func Info(values ...any) {
	writeMessage(Default().packageEvent(zerolog.InfoLevel), values)
}

// Warn writes a warning message through the default Logger.
func Warn(values ...any) {
	writeMessage(Default().packageEvent(zerolog.WarnLevel), values)
}

// Error writes an error message through the default Logger.
func Error(values ...any) {
	writeMessage(Default().packageEvent(zerolog.ErrorLevel), values)
}

// Fatal writes a fatal message and terminates the process with exit status 1
// when fatal logging is enabled.
func Fatal(values ...any) {
	writeMessage(Default().packageEvent(zerolog.FatalLevel), values)
}

// Panic writes a panic message and then panics when panic logging is enabled.
func Panic(values ...any) {
	writeMessage(Default().packageEvent(zerolog.PanicLevel), values)
}

func (logger *Logger) packageEvent(level zerolog.Level) *zerolog.Event {
	if !logger.level.enabled(level) {
		return nil
	}
	switch level {
	case zerolog.TraceLevel:
		return logger.packageLevelLogger.Trace()
	case zerolog.DebugLevel:
		return logger.packageLevelLogger.Debug()
	case zerolog.InfoLevel:
		return logger.packageLevelLogger.Info()
	case zerolog.WarnLevel:
		return logger.packageLevelLogger.Warn()
	case zerolog.ErrorLevel:
		return logger.packageLevelLogger.Error()
	case zerolog.FatalLevel:
		return logger.packageLevelLogger.Fatal()
	case zerolog.PanicLevel:
		return logger.packageLevelLogger.Panic()
	default:
		return nil
	}
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
