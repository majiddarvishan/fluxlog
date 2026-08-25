package fluxlog

import (
	"errors"
	"io"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog"
)

// managedWriter serializes writes and makes Close safe to call while other
// goroutines are logging.
type managedWriter struct {
	mu      sync.Mutex
	writer  zerolog.LevelWriter
	closers []io.Closer
	closed  bool
}

func (writer *managedWriter) Write(data []byte) (int, error) {
	return writer.WriteLevel(zerolog.NoLevel, data)
}

func (writer *managedWriter) WriteLevel(level zerolog.Level, data []byte) (int, error) {
	writer.mu.Lock()
	defer writer.mu.Unlock()

	if writer.closed {
		return 0, io.ErrClosedPipe
	}
	return writer.writer.WriteLevel(level, data)
}

func (writer *managedWriter) Close() error {
	writer.mu.Lock()
	defer writer.mu.Unlock()

	if writer.closed {
		return nil
	}
	writer.closed = true

	var closeErrors []error
	for index := len(writer.closers) - 1; index >= 0; index-- {
		if err := writer.closers[index].Close(); err != nil {
			closeErrors = append(closeErrors, err)
		}
	}
	return errors.Join(closeErrors...)
}

// levelWriter provides per-instance runtime level changes. It deliberately
// avoids zerolog.SetGlobalLevel, so independent Logger instances cannot affect
// one another.
type levelWriter struct {
	output *managedWriter
	level  atomic.Int32
}

func newLevelWriter(output *managedWriter, level zerolog.Level) *levelWriter {
	writer := &levelWriter{output: output}
	writer.level.Store(int32(level))
	return writer
}

func (writer *levelWriter) Write(data []byte) (int, error) {
	return writer.WriteLevel(zerolog.NoLevel, data)
}

func (writer *levelWriter) WriteLevel(level zerolog.Level, data []byte) (int, error) {
	if !writer.enabled(level) {
		return len(data), nil
	}
	return writer.output.WriteLevel(level, data)
}

func (writer *levelWriter) enabled(level zerolog.Level) bool {
	minimum := zerolog.Level(writer.level.Load())
	return minimum != zerolog.Disabled && (level == zerolog.NoLevel || level >= minimum)
}

func (writer *levelWriter) setLevel(level zerolog.Level) {
	writer.level.Store(int32(level))
}

func (writer *levelWriter) currentLevel() zerolog.Level {
	return zerolog.Level(writer.level.Load())
}

func (writer *levelWriter) Close() error {
	return writer.output.Close()
}
