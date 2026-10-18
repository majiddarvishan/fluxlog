package fluxlog

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// LVL-01: zerolog.GlobalLevel filters before the fluxlog per-instance writer.
// Restore global state: this test must never call t.Parallel.
func TestLVL01ExternalZerologGlobalThreshold(t *testing.T) {
	original := zerolog.GlobalLevel()
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	t.Cleanup(func() { zerolog.SetGlobalLevel(original) })

	var output bytes.Buffer
	logger, err := New(Config{
		Level: DebugLevel,
		Console: &ConsoleConfig{
			Writer: &output,
			Format: JSONFormat,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()
	logger.Debug().Msg("should-be-filtered")
	logger.Error().Msg("should-be-visible")
	if strings.Contains(output.String(), "should-be-filtered") {
		t.Fatalf("external zerolog global threshold did not filter debug: %q", output.String())
	}
	if !strings.Contains(output.String(), "should-be-visible") {
		t.Fatalf("error log did not survive global threshold: %q", output.String())
	}
	if got := logger.Level(); got != DebugLevel {
		t.Fatalf("fluxlog instance threshold was unexpectedly changed to %q", got)
	}
}
