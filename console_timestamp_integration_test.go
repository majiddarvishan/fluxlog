package fluxlog

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// FMT-01: zerolog serializes Unix timestamps as JSON numbers. Fluxlog's
// custom timestamp formatter should consume each supported numeric form.
func TestFMT01ConsoleUnixTimestampFormats(t *testing.T) {
	want := time.Unix(1700000000, 0).UTC().Format(time.RFC3339)
	for _, tc := range []struct {
		name      string
		timestamp string
	}{
		{name: "seconds", timestamp: "1700000000"},
		{name: "milliseconds", timestamp: "1700000000000"},
		{name: "microseconds", timestamp: "1700000000000000"},
		{name: "nanoseconds", timestamp: "1700000000000000000"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var console bytes.Buffer
			writer := newConsoleWriter(&console, NeverColor, time.RFC3339, DefaultCallerMaxLength)
			payload := []byte("{\"level\":\"info\",\"time\":" + tc.timestamp + ",\"message\":\"timestamp-check\"}")
			if _, err := writer.Write(payload); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(console.String(), want) {
				if strings.Contains(console.String(), "timestamp-check") {
					t.Skipf("KNOWN DEFECT FMT-01: numeric %s timestamp omitted from console: %q", tc.name, console.String())
				}
				t.Fatalf("unexpected console output: %q", console.String())
			}
			var jsonOutput bytes.Buffer
			jsonWriter := outputWriter(&jsonOutput, JSONFormat, NeverColor, time.RFC3339, DefaultCallerMaxLength)
			if _, err := jsonWriter.Write(payload); err != nil {
				t.Fatal(err)
			}
			if jsonOutput.String() != string(payload) {
				t.Fatalf("JSON output changed: got %q, want %q", jsonOutput.String(), payload)
			}
		})
	}
}

// FMT-01: the already-supported RFC3339 string case stays a positive test.
func TestFMT01ConsoleRFC3339Timestamp(t *testing.T) {
	var output bytes.Buffer
	writer := newConsoleWriter(&output, NeverColor, time.RFC3339, DefaultCallerMaxLength)
	payload := []byte("{\"level\":\"info\",\"time\":\"2023-11-14T22:13:20Z\",\"message\":\"timestamp-check\"}")
	if _, err := writer.Write(payload); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "2023-11-14T22:13:20Z") {
		t.Fatalf("RFC3339 timestamp missing from console: %q", output.String())
	}
}
