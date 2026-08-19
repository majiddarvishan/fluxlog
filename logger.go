// Package fluxlog provides instance-based structured logging with runtime-safe
// level changes and optional console and rotating-file outputs.
package fluxlog

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger is an independent logger instance. It is safe for concurrent use.
// Call Close when a file output is configured.
type Logger struct {
	logger    zerolog.Logger
	level     *levelWriter
	output    *managedWriter
	closeOnce sync.Once
	closeErr  error
}

// New validates config and constructs an independent Logger without changing
// zerolog's process-global logger or global level.
func New(config Config) (*Logger, error) {
	normalized, err := config.normalized()
	if err != nil {
		return nil, err
	}

	minimumLevel, err := normalized.Level.zerologLevel()
	if err != nil {
		return nil, err
	}

	writers := make([]io.Writer, 0, 2)
	closers := make([]io.Closer, 0, 1)

	if normalized.Console != nil {
		writers = append(writers, outputWriter(
			normalized.Console.Writer,
			normalized.Console.Format,
			normalized.Console.Color,
			normalized.Console.TimeFormat,
		))
	}

	if normalized.File != nil {
		if err := os.MkdirAll(filepath.Dir(normalized.File.Path), 0o750); err != nil {
			return nil, fmt.Errorf("fluxlog: create log directory: %w", err)
		}
		fileOutput := &lumberjack.Logger{
			Filename:   normalized.File.Path,
			MaxSize:    normalized.File.MaxSizeMB,
			MaxBackups: normalized.File.MaxBackups,
			MaxAge:     normalized.File.MaxAgeDays,
			Compress:   normalized.File.Compress,
			LocalTime:  normalized.File.LocalTime,
		}
		writers = append(writers, outputWriter(
			fileOutput,
			normalized.File.Format,
			NeverColor,
			normalizedConsoleTimeFormat(normalized.Console),
		))
		closers = append(closers, fileOutput)
	}

	multi := zerolog.MultiLevelWriter(writers...)
	managed := &managedWriter{writer: multi, closers: closers}
	dynamicLevel := newLevelWriter(managed, minimumLevel)

	base := zerolog.New(dynamicLevel)
	contextBuilder := base.With()
	if normalized.Timestamp {
		contextBuilder = contextBuilder.Timestamp()
	}
	if normalized.Service != "" {
		contextBuilder = contextBuilder.Str("service", normalized.Service)
	}
	if normalized.Caller {
		contextBuilder = contextBuilder.Caller()
	}

	return &Logger{
		logger: contextBuilder.Logger(),
		level:  dynamicLevel,
		output: managed,
	}, nil
}

func outputWriter(output io.Writer, format OutputFormat, color ColorMode, timeFormat string) io.Writer {
	if format == JSONFormat {
		return output
	}

	return zerolog.ConsoleWriter{
		Out:        output,
		NoColor:    !colorEnabled(color, output),
		TimeFormat: timeFormat,
	}
}

func normalizedConsoleTimeFormat(config *ConsoleConfig) string {
	if config == nil {
		return DefaultConfig().Console.TimeFormat
	}
	return config.TimeFormat
}

func colorEnabled(mode ColorMode, output io.Writer) bool {
	switch mode {
	case AlwaysColor:
		return true
	case NeverColor:
		return false
	case AutoColor:
		file, ok := output.(*os.File)
		if !ok {
			return false
		}
		info, err := file.Stat()
		return err == nil && info.Mode()&os.ModeCharDevice != 0
	default:
		return false
	}
}

// SetLevel changes only this Logger's minimum level and is safe to call while
// other goroutines are writing logs.
func (logger *Logger) SetLevel(level Level) error {
	zerologLevel, err := level.zerologLevel()
	if err != nil {
		return err
	}
	logger.level.setLevel(zerologLevel)
	return nil
}

// SetLevelString parses and applies a level name.
func (logger *Logger) SetLevelString(level string) error {
	parsed, err := ParseLevel(level)
	if err != nil {
		return err
	}
	return logger.SetLevel(parsed)
}

// Level returns this Logger's current minimum severity.
func (logger *Logger) Level() Level {
	switch logger.level.currentLevel() {
	case zerolog.TraceLevel:
		return TraceLevel
	case zerolog.DebugLevel:
		return DebugLevel
	case zerolog.InfoLevel:
		return InfoLevel
	case zerolog.WarnLevel:
		return WarnLevel
	case zerolog.ErrorLevel:
		return ErrorLevel
	case zerolog.FatalLevel:
		return FatalLevel
	case zerolog.PanicLevel:
		return PanicLevel
	case zerolog.Disabled:
		return DisabledLevel
	default:
		return DisabledLevel
	}
}

// Zerolog returns a copy of the underlying zerolog.Logger for advanced usage.
// The copy still uses fluxlog's runtime level and managed output lifecycle.
func (logger *Logger) Zerolog() zerolog.Logger {
	return logger.logger
}

// With returns a zerolog context derived from this Logger.
func (logger *Logger) With() zerolog.Context {
	return logger.logger.With()
}

// WithContext stores this Logger in ctx for retrieval with zerolog.Ctx.
func (logger *Logger) WithContext(ctx context.Context) context.Context {
	return logger.logger.WithContext(ctx)
}

func (logger *Logger) Trace() *zerolog.Event {
	if !logger.level.enabled(zerolog.TraceLevel) {
		return nil
	}
	return logger.logger.Trace()
}

func (logger *Logger) Debug() *zerolog.Event {
	if !logger.level.enabled(zerolog.DebugLevel) {
		return nil
	}
	return logger.logger.Debug()
}

func (logger *Logger) Info() *zerolog.Event {
	if !logger.level.enabled(zerolog.InfoLevel) {
		return nil
	}
	return logger.logger.Info()
}

func (logger *Logger) Warn() *zerolog.Event {
	if !logger.level.enabled(zerolog.WarnLevel) {
		return nil
	}
	return logger.logger.Warn()
}

func (logger *Logger) Error() *zerolog.Event {
	if !logger.level.enabled(zerolog.ErrorLevel) {
		return nil
	}
	return logger.logger.Error()
}

// Fatal returns a fatal event. Calling Msg on it terminates the process. Core
// fluxlog internals never use Fatal; this method only exposes zerolog behavior
// when the application explicitly requests it.
func (logger *Logger) Fatal() *zerolog.Event {
	if !logger.level.enabled(zerolog.FatalLevel) {
		return nil
	}
	return logger.logger.Fatal()
}

// Panic returns a panic event. Calling Msg on it panics after writing.
func (logger *Logger) Panic() *zerolog.Event {
	if !logger.level.enabled(zerolog.PanicLevel) {
		return nil
	}
	return logger.logger.Panic()
}

// Close releases owned file outputs. It is safe to call more than once.
func (logger *Logger) Close() error {
	logger.closeOnce.Do(func() {
		logger.closeErr = logger.output.Close()
	})
	return logger.closeErr
}
