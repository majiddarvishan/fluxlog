package fluxlog

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestConsolePattern(t *testing.T) {
	var output bytes.Buffer
	writer := newConsoleWriter(&output, NeverColor, time.RFC3339, DefaultCallerMaxLength)
	event := []byte(`{"level":"info","time":"2026-08-25T18:22:43+03:30","caller":"internal/webservice/web_server.go:272","service":"gateway","message":"Web server starting on 0.0.0.0:4567"}`)

	if _, err := writer.Write(event); err != nil {
		t.Fatal(err)
	}

	want := "2026-08-25T18:22:43+03:30 INF [ …_server.go:272 ] (Gateway) Web server starting on 0.0.0.0:4567\n"
	if output.String() != want {
		t.Fatalf("console output = %q, want %q", output.String(), want)
	}
}

func TestConsoleServiceUsesColor(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	var output bytes.Buffer
	writer := newConsoleWriter(&output, AlwaysColor, time.RFC3339, DefaultCallerMaxLength)
	event := []byte(`{"level":"info","service":"gateway","message":"started"}`)

	if _, err := writer.Write(event); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), ansiCyan+"(Gateway)"+ansiReset) {
		t.Fatalf("colored service missing from output: %q", output.String())
	}
}

func TestConsoleServiceRespectsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var output bytes.Buffer
	writer := newConsoleWriter(&output, AlwaysColor, time.RFC3339, DefaultCallerMaxLength)
	event := []byte(`{"level":"info","service":"gateway","message":"started"}`)

	if _, err := writer.Write(event); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "\x1b[") {
		t.Fatalf("NO_COLOR output contains ANSI codes: %q", output.String())
	}
}

func TestConsolePatternKeepsAdditionalFields(t *testing.T) {
	var output bytes.Buffer
	writer := newConsoleWriter(&output, NeverColor, time.RFC3339, DefaultCallerMaxLength)
	event := []byte(`{"level":"info","message":"received","service":"gateway","request_id":"req-1"}`)

	if _, err := writer.Write(event); err != nil {
		t.Fatal(err)
	}

	want := "INF (Gateway) received request_id=req-1\n"
	if output.String() != want {
		t.Fatalf("console output = %q, want %q", output.String(), want)
	}
}

func TestTruncateLeftUsesRuneLength(t *testing.T) {
	if actual := truncateLeft("مسیر/فایل.go:10", 8); actual != "…ل.go:10" {
		t.Fatalf("truncateLeft returned %q", actual)
	}
}
