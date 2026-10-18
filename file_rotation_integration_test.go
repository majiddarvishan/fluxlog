package fluxlog

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/natefinch/lumberjack.v2"
)

// FILE-01: two individually valid messages crossing MaxSize must rotate on
// disk. Both active and rotated content are checked without a timing delay.
func TestFILE01SizeRotationPreservesRecords(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "application.log")
	logger, err := New(Config{
		Level: InfoLevel,
		File:  &FileConfig{Path: path, Format: JSONFormat, MaxSizeMB: 1, MaxBackups: 3},
	})
	if err != nil {
		t.Fatal(err)
	}
	payload := strings.Repeat("x", 620_000)
	logger.Info().Str("payload", payload).Msg("rotation-first")
	logger.Info().Str("payload", payload).Msg("rotation-second")
	if err := logger.Close(); err != nil {
		t.Fatal(err)
	}

	current, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(current, []byte("rotation-second")) || bytes.Contains(current, []byte("rotation-first")) {
		t.Fatal("the active log file does not contain exactly the post-rotation record")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	backups := 0
	for _, entry := range entries {
		if entry.Name() == "application.log" || entry.IsDir() {
			continue
		}
		if !strings.HasPrefix(entry.Name(), "application-") || !strings.HasSuffix(entry.Name(), ".log") {
			t.Fatalf("unexpected file in isolated log directory: %q", entry.Name())
		}
		info, err := entry.Info()
		if err != nil {
			t.Fatal(err)
		}
		if info.Size() > 1<<20 {
			t.Fatalf("backup %q exceeds 1MiB: %d", entry.Name(), info.Size())
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(data, []byte("rotation-first")) {
			t.Fatalf("backup %q is missing the original record", entry.Name())
		}
		backups++
	}
	if backups != 1 {
		t.Fatalf("backup count = %d, want 1", backups)
	}
	if int64(len(current)) > 1<<20 {
		t.Fatalf("active log file exceeds 1MiB: %d", len(current))
	}
}

// FILE-01: retention and naming settings are passed unchanged to lumberjack.
// Async compression and stale-backup cleanup need separate execution tests.
func TestFILE01RotationConfigurationIsForwarded(t *testing.T) {
	logger, err := New(Config{
		Level: InfoLevel,
		File: &FileConfig{
			Path:       filepath.Join(t.TempDir(), "audit.log"),
			Format:     JSONFormat,
			MaxSizeMB:  2,
			MaxBackups: 2,
			MaxAgeDays: 7,
			Compress:   true,
			LocalTime:  true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()
	if len(logger.output.closers) != 1 {
		t.Fatalf("owned closers = %d, want 1", len(logger.output.closers))
	}
	rotator, ok := logger.output.closers[0].(*lumberjack.Logger)
	if !ok {
		t.Fatalf("owned file closer is %T, want *lumberjack.Logger", logger.output.closers[0])
	}
	if rotator.MaxSize != 2 || rotator.MaxBackups != 2 || rotator.MaxAge != 7 || !rotator.Compress || !rotator.LocalTime {
		t.Fatalf("rotation options were not forwarded: %+v", rotator)
	}
}
