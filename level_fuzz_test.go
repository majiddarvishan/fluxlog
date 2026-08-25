package fluxlog

import "testing"

func FuzzParseLevel(f *testing.F) {
	for _, seed := range []string{
		"trace", "debug", "info", "warn", "warning", "error",
		"fatal", "panic", "disabled", "off", " DEBUG ", "unknown",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, value string) {
		level, err := ParseLevel(value)
		if err != nil {
			return
		}
		if _, err := level.zerologLevel(); err != nil {
			t.Fatalf("ParseLevel returned an unusable level %q: %v", level, err)
		}
	})
}
