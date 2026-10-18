package fluxlog

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

// FMT-01: zerolog serializes Unix timestamps as JSON numbers. Fluxlog's
// custom timestamp formatter should consume each supported numeric form.
func TestFMT01ConsoleUnixTimestampFormats(t *testing.T) {
	want := time.Unix(1700000000, 0).Local().Format(time.RFC3339)
	for _, tc := range []struct {
		name      string
		timestamp string
		format    string
	}{
		{name: "seconds", timestamp: "1700000000", format: zerolog.TimeFormatUnix},
		{name: "milliseconds", timestamp: "1700000000000", format: zerolog.TimeFormatUnixMs},
		{name: "microseconds", timestamp: "1700000000000000", format: zerolog.TimeFormatUnixMicro},
		{name: "nanoseconds", timestamp: "1700000000000000000", format: zerolog.TimeFormatUnixNano},
	} {
		t.Run(tc.name, func(t *testing.T) {
			previous := zerolog.TimeFieldFormat
			zerolog.TimeFieldFormat = tc.format
			t.Cleanup(func() { zerolog.TimeFieldFormat = previous })

			var console bytes.Buffer
			writer := newConsoleWriter(&console, NeverColor, time.RFC3339, DefaultCallerMaxLength)
			payload := []byte("{\"level\":\"info\",\"time\":" + tc.timestamp + ",\"message\":\"timestamp-check\"}")
			if _, err := writer.Write(payload); err != nil {
				t.Fatal(err)
			}
			var jsonOutput bytes.Buffer
			jsonWriter := outputWriter(&jsonOutput, JSONFormat, NeverColor, time.RFC3339, DefaultCallerMaxLength)
			if _, err := jsonWriter.Write(payload); err != nil {
				t.Fatal(err)
			}
			if jsonOutput.String() != string(payload) {
				t.Fatalf("JSON output changed: got %q, want %q", jsonOutput.String(), payload)
			}
			if !strings.Contains(console.String(), want) {
				t.Fatalf("human-readable %s timestamp missing from console: %q, want %q", tc.name, console.String(), want)
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


// FMT-01: use the application-configured layout for an integer timestamp,
// rather than exposing the raw epoch value in human-readable output.
func TestFMT01CustomHumanReadableTimeLayout(t *testing.T) {
	previous := zerolog.TimeFieldFormat
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	t.Cleanup(func() { zerolog.TimeFieldFormat = previous })

	const layout = "2006-01-02 15:04:05"
	var output bytes.Buffer
	writer := newConsoleWriter(&output, NeverColor, layout, DefaultCallerMaxLength)
	payload := []byte(`{"level":"info","time":1700000000000,"message":"timestamp-check"}`)
	if _, err := writer.Write(payload); err != nil {
		t.Fatal(err)
	}
	want := time.Unix(1700000000, 0).Local().Format(layout) + " INF timestamp-check\n"
	if got := output.String(); got != want {
		t.Fatalf("console = %q, want %q", got, want)
	}
}

// Fractional timestamps must retain precision, and negative epoch values
// must be converted without overflowing a time.Duration.
func TestFMT01UnixTimestampPrecisionAndNegativeEpoch(t *testing.T) {
	for _, tc := range []struct {
		name      string
		format    string
		value     string
		expected  time.Time
	}{
		{name: "millisecond-fraction", format: zerolog.TimeFormatUnixMs, value: "1700000000123", expected: time.Unix(1700000000, 123_000_000)},
		{name: "microsecond-fraction", format: zerolog.TimeFormatUnixMicro, value: "1700000000123456", expected: time.Unix(1700000000, 123_456_000)},
		{name: "nanosecond-fraction", format: zerolog.TimeFormatUnixNano, value: "1700000000123456789", expected: time.Unix(1700000000, 123_456_789)},
		{name: "negative-milliseconds", format: zerolog.TimeFormatUnixMs, value: "-1", expected: time.Unix(-1, 999_000_000)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			previous := zerolog.TimeFieldFormat
			zerolog.TimeFieldFormat = tc.format
			t.Cleanup(func() { zerolog.TimeFieldFormat = previous })

			var output bytes.Buffer
			writer := newConsoleWriter(&output, NeverColor, time.RFC3339Nano, DefaultCallerMaxLength)
			payload := []byte(`{"level":"info","time":` + tc.value + `,"message":"precision-check"}`)
			if _, err := writer.Write(payload); err != nil {
				t.Fatal(err)
			}
			want := tc.expected.Local().Format(time.RFC3339Nano) + " INF precision-check\n"
			if got := output.String(); got != want {
				t.Fatalf("console = %q, want %q", got, want)
			}
		})
	}
}

// Verify end-to-end output from a real Logger with numeric zerolog timestamps.
func TestFMT01LoggerRendersUnixTimestampAsHumanReadable(t *testing.T) {
	previousFormat := zerolog.TimeFieldFormat
	previousClock := zerolog.TimestampFunc
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano
	zerolog.TimestampFunc = func() time.Time { return time.Unix(1700000000, 123_456_789).UTC() }
	t.Cleanup(func() {
		zerolog.TimeFieldFormat = previousFormat
		zerolog.TimestampFunc = previousClock
	})

	const layout = "2006-01-02 15:04:05.000"
	var output bytes.Buffer
	logger, err := New(Config{
		Level: InfoLevel,
		Timestamp: true,
		Console: &ConsoleConfig{
			Writer: &output,
			Format: ConsoleFormat,
			Color: NeverColor,
			TimeFormat: layout,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()
	logger.Info().Msg("human-time")
	want := time.Unix(1700000000, 123_456_789).Local().Format(layout) + " INF human-time\n"
	if got := output.String(); got != want {
		t.Fatalf("logger output = %q, want %q", got, want)
	}
}
