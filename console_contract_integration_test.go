package fluxlog

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// FMT-02: color is opt-in for a non-terminal buffer, and explicit modes
// do not accidentally add or remove ANSI sequences.
func TestFMT02ConsoleColorModes(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	for _, tc := range []struct {
		name      string
		color     ColorMode
		expectANSI bool
	}{
		{name: "auto-buffer", color: AutoColor, expectANSI: false},
		{name: "always", color: AlwaysColor, expectANSI: true},
		{name: "never", color: NeverColor, expectANSI: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			writer := newConsoleWriter(&output, tc.color, time.RFC3339, DefaultCallerMaxLength)
			if _, err := writer.Write([]byte(`{"level":"info","service":"gateway","message":"started"}`)); err != nil {
				t.Fatal(err)
			}
			rendered := output.String()
			if got := strings.Contains(rendered, "\x1b["); got != tc.expectANSI {
				t.Fatalf("color mode %s: ANSI present = %t, want %t: %q", tc.color, got, tc.expectANSI, rendered)
			}
			if !strings.Contains(rendered, "(Gateway)") || !strings.Contains(rendered, "started") {
				t.Fatalf("service/message missing from console output: %q", rendered)
			}
		})
	}
}

// FMT-02: NO_COLOR wins even when AlwaysColor was explicitly chosen.
func TestFMT02NoColorDisablesServiceLevelAndCallerColors(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var output bytes.Buffer
	writer := newConsoleWriter(&output, AlwaysColor, time.RFC3339, DefaultCallerMaxLength)
	event := []byte(`{"level":"info","caller":"app/handler.go:123","service":"gateway","message":"received","request_id":"req-1"}`)
	if _, err := writer.Write(event); err != nil {
		t.Fatal(err)
	}
	rendered := output.String()
	if strings.Contains(rendered, "\x1b[") {
		t.Fatalf("NO_COLOR console output still has ANSI codes: %q", rendered)
	}
	for _, field := range []string{"INF", "(Gateway)", "received", "request_id=req-1", "handler.go:123"} {
		if !strings.Contains(rendered, field) {
			t.Fatalf("NO_COLOR output is missing %q: %q", field, rendered)
		}
	}
}

// FMT-02: a ConsoleFormat file sink is always plain text, even if a
// simultaneously configured console stream forces ANSI colors.
func TestFMT02DualOutputKeepsFilePlainText(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	path := filepath.Join(t.TempDir(), "human.log")
	var console bytes.Buffer
	logger, err := New(Config{
		Level: InfoLevel,
		Service: "gateway",
		Caller: true,
		Console: &ConsoleConfig{
			Writer: &console,
			Format: ConsoleFormat,
			Color: AlwaysColor,
		},
		File: &FileConfig{
			Path: path,
			Format: ConsoleFormat,
			MaxSizeMB: 1,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	logger.Info().Str("request_id", "req-2").Msg("dual-output")
	if err := logger.Close(); err != nil {
		t.Fatal(err)
	}
	streamOutput := console.String()
	if !strings.Contains(streamOutput, "\x1b[") || !strings.Contains(streamOutput, "(Gateway)") {
		t.Fatalf("forced-color console missing expected formatting: %q", streamOutput)
	}
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	fileOutput := string(fileBytes)
	if strings.Contains(fileOutput, "\x1b[") {
		t.Fatalf("file console format unexpectedly has ANSI: %q", fileOutput)
	}
	for _, field := range []string{"INF", "(Gateway)", "dual-output", "request_id=req-2", ".go:"} {
		if !strings.Contains(fileOutput, field) {
			t.Fatalf("file output missing %q: %q", field, fileOutput)
		}
	}
}
