package fluxlog

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"testing"
)

func newTestLogger(t *testing.T, level Level, output *bytes.Buffer) *Logger {
	t.Helper()

	logger, err := New(Config{
		Level:     level,
		Timestamp: false,
		Console: &ConsoleConfig{
			Writer: output,
			Format: JSONFormat,
		},
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	t.Cleanup(func() {
		if err := logger.Close(); err != nil {
			t.Errorf("Close returned error: %v", err)
		}
	})
	return logger
}

func TestStructuredLog(t *testing.T) {
	var output bytes.Buffer
	logger := newTestLogger(t, InfoLevel, &output)

	logger.Info().Str("request_id", "req-1").Msg("request received")

	var event map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(output.Bytes()), &event); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if event["level"] != "info" {
		t.Fatalf("level = %v, want info", event["level"])
	}
	if event["request_id"] != "req-1" {
		t.Fatalf("request_id = %v, want req-1", event["request_id"])
	}
	if event["message"] != "request received" {
		t.Fatalf("message = %v, want request received", event["message"])
	}
}

func TestRuntimeLevelIsPerInstance(t *testing.T) {
	var firstOutput bytes.Buffer
	var secondOutput bytes.Buffer
	first := newTestLogger(t, InfoLevel, &firstOutput)
	second := newTestLogger(t, InfoLevel, &secondOutput)

	if err := first.SetLevel(DebugLevel); err != nil {
		t.Fatalf("SetLevel returned error: %v", err)
	}
	first.Debug().Msg("first")
	second.Debug().Msg("second")

	if !strings.Contains(firstOutput.String(), "first") {
		t.Fatal("first logger did not apply its runtime level")
	}
	if secondOutput.Len() != 0 {
		t.Fatal("changing first logger affected second logger")
	}
}

func TestConcurrentLoggingAndLevelChanges(t *testing.T) {
	var output bytes.Buffer
	logger := newTestLogger(t, InfoLevel, &output)

	var workers sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		workers.Add(1)
		go func(id int) {
			defer workers.Done()
			for item := 0; item < 100; item++ {
				logger.Info().Int("worker", id).Int("item", item).Msg("event")
				if item%10 == 0 {
					_ = logger.SetLevel(DebugLevel)
					_ = logger.SetLevel(InfoLevel)
				}
			}
		}(worker)
	}
	workers.Wait()

	scanner := bufio.NewScanner(bytes.NewReader(output.Bytes()))
	count := 0
	for scanner.Scan() {
		var event map[string]any
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			t.Fatalf("concurrent output contains invalid JSON: %v", err)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("reading output: %v", err)
	}
	if count != 800 {
		t.Fatalf("event count = %d, want 800", count)
	}
}

func TestCloseIsIdempotent(t *testing.T) {
	var output bytes.Buffer
	logger := newTestLogger(t, InfoLevel, &output)

	if err := logger.Close(); err != nil {
		t.Fatalf("first Close returned error: %v", err)
	}
	if err := logger.Close(); err != nil {
		t.Fatalf("second Close returned error: %v", err)
	}
}

func TestDisabledEventDoesNotEvaluateFields(t *testing.T) {
	var output bytes.Buffer
	logger := newTestLogger(t, InfoLevel, &output)

	event := logger.Debug()
	if event != nil {
		t.Fatal("expected a nil zerolog event for a disabled level")
	}
	event.Str("expensive", "value").Msg("disabled")
	if output.Len() != 0 {
		t.Fatal("disabled event was written")
	}
}

func TestFileOutputCreatesParentDirectory(t *testing.T) {
	path := t.TempDir() + "/nested/application.log"
	logger, err := New(Config{
		Level: InfoLevel,
		File:  &FileConfig{Path: path, Format: JSONFormat},
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	logger.Info().Msg("written to file")
	if err := logger.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if !bytes.Contains(content, []byte("written to file")) {
		t.Fatalf("file does not contain the log event: %s", content)
	}
}
