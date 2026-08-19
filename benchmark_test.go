package fluxlog

import (
	"io"
	"testing"
)

func benchmarkLogger(b *testing.B, level Level) *Logger {
	b.Helper()
	logger, err := New(Config{
		Level:   level,
		Console: &ConsoleConfig{
			Writer: io.Discard,
			Format: JSONFormat,
		},
	})
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() { _ = logger.Close() })
	return logger
}

func BenchmarkEnabledLog(b *testing.B) {
	logger := benchmarkLogger(b, InfoLevel)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info().Int("iteration", i).Msg("benchmark")
	}
}

func BenchmarkDisabledLog(b *testing.B) {
	logger := benchmarkLogger(b, InfoLevel)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Debug().Int("iteration", i).Msg("benchmark")
	}
}
