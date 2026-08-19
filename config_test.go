package fluxlog

import (
	"bytes"
	"testing"
)

func TestConfigRequiresOutput(t *testing.T) {
	_, err := New(Config{Level: InfoLevel})
	if err == nil {
		t.Fatal("expected an error when no output is configured")
	}
}

func TestConfigRejectsInvalidLevel(t *testing.T) {
	_, err := New(Config{
		Level:   Level("verbose"),
		Console: &ConsoleConfig{Writer: &bytes.Buffer{}, Format: JSONFormat},
	})
	if err == nil {
		t.Fatal("expected invalid level to be rejected")
	}
}

func TestConfigRejectsNegativeRotationValue(t *testing.T) {
	_, err := New(Config{
		Level: InfoLevel,
		File:  &FileConfig{Path: "test.log", MaxBackups: -1},
	})
	if err == nil {
		t.Fatal("expected negative rotation value to be rejected")
	}
}

func TestParseLevelAliases(t *testing.T) {
	tests := map[string]Level{
		"TRACE":   TraceLevel,
		" warning ": WarnLevel,
		"off":     DisabledLevel,
	}

	for input, expected := range tests {
		actual, err := ParseLevel(input)
		if err != nil {
			t.Fatalf("ParseLevel(%q) returned error: %v", input, err)
		}
		if actual != expected {
			t.Fatalf("ParseLevel(%q) = %q, want %q", input, actual, expected)
		}
	}
}
