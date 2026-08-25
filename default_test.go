package fluxlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

type panicStringer struct{}

func (panicStringer) String() string {
	panic("disabled log arguments must not be formatted")
}

func TestPackageLevelLogging(t *testing.T) {
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

	previous := SetDefault(logger)
	defer SetDefault(previous)

	Info("server started")
	Debug("processed %d requests", 3)
	Error(errors.New("connection failed"))

	logs := output.String()
	if !strings.Contains(logs, `"message":"server started"`) {
		t.Fatalf("plain message missing from output: %s", logs)
	}
	if !strings.Contains(logs, `"message":"processed 3 requests"`) {
		t.Fatalf("formatted message missing from output: %s", logs)
	}
	if !strings.Contains(logs, `"message":"connection failed"`) {
		t.Fatalf("non-string message missing from output: %s", logs)
	}
}

func TestPackageLevelLoggingHonorsRuntimeLevel(t *testing.T) {
	var output bytes.Buffer
	logger, err := New(Config{
		Level: InfoLevel,
		Console: &ConsoleConfig{
			Writer: &output,
			Format: JSONFormat,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	previous := SetDefault(logger)
	defer SetDefault(previous)

	Debug("hidden %s", panicStringer{})
	if output.Len() != 0 {
		t.Fatalf("disabled debug message was written: %s", output.String())
	}

	if err := logger.SetLevel(DebugLevel); err != nil {
		t.Fatal(err)
	}
	Debug("visible %d", 2)
	if !strings.Contains(output.String(), `"message":"visible 2"`) {
		t.Fatalf("runtime level change was not applied: %s", output.String())
	}
}

func TestSetDefaultRejectsNil(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("SetDefault(nil) did not panic")
		}
	}()
	SetDefault(nil)
}

func TestPackageLevelCallerPointsToApplication(t *testing.T) {
	var output bytes.Buffer
	logger, err := New(Config{
		Level:  InfoLevel,
		Caller: true,
		Console: &ConsoleConfig{
			Writer: &output,
			Format: JSONFormat,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()

	previous := SetDefault(logger)
	defer SetDefault(previous)
	Info("caller test")

	var event map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(output.Bytes()), &event); err != nil {
		t.Fatal(err)
	}
	caller, _ := event["caller"].(string)
	if !strings.Contains(caller, "default_test.go:") {
		t.Fatalf("caller = %q, want application call site", caller)
	}
}
