package fluxlog

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// The lumberjack retention mill is asynchronous. Poll its externally
// observable result with a bounded deadline rather than sleeping blindly.
func waitForRotationCleanup(t *testing.T, description string, condition func() bool) {
	t.Helper()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()
	for {
		if condition() {
			return
		}
		select {
		case <-ticker.C:
		case <-timeout.C:
			t.Fatalf("timed out waiting for %s", description)
		}
	}
}

func seededBackupName(base string, at time.Time) string {
	return base + "-" + at.UTC().Format("2006-01-02T15-04-05.000") + ".log"
}

// FILE-01: exercising the real asynchronous mill with seeded dated backups
// verifies that MaxBackups removes older entries, not just config forwarding.
func TestFILE01MaxBackupsIsEnforced(t *testing.T) {
	dir := t.TempDir()
	now := time.Now().UTC()
	old := filepath.Join(dir, seededBackupName("service", now.Add(-3*time.Hour)))
	middle := filepath.Join(dir, seededBackupName("service", now.Add(-2*time.Hour)))
	recent := filepath.Join(dir, seededBackupName("service", now.Add(-time.Hour)))
	for _, name := range []string{old, middle, recent} {
		if err := os.WriteFile(name, []byte("historical"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	logger, err := New(Config{
		Level: InfoLevel,
		File: &FileConfig{
			Path:       filepath.Join(dir, "service.log"),
			Format:     JSONFormat,
			MaxSizeMB:  1,
			MaxBackups: 1,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()
	logger.Info().Msg("start retention")
	waitForRotationCleanup(t, "MaxBackups removal", func() bool {
		_, oldErr := os.Stat(old)
		_, middleErr := os.Stat(middle)
		_, recentErr := os.Stat(recent)
		return os.IsNotExist(oldErr) && os.IsNotExist(middleErr) && recentErr == nil
	})
}

// FILE-01: verify age expiration and actual gzip archive contents, not merely
// that Compress and MaxAge are copied to the lumberjack config object.
func TestFILE01MaxAgeAndCompression(t *testing.T) {
	dir := t.TempDir()
	now := time.Now().UTC()
	stale := filepath.Join(dir, seededBackupName("service", now.AddDate(0, 0, -30)))
	recent := filepath.Join(dir, seededBackupName("service", now.Add(-time.Hour)))
	if err := os.WriteFile(stale, []byte("old-data"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(recent, []byte("recent-data"), 0o600); err != nil {
		t.Fatal(err)
	}
	logger, err := New(Config{
		Level: InfoLevel,
		File: &FileConfig{
			Path:       filepath.Join(dir, "service.log"),
			Format:     JSONFormat,
			MaxSizeMB:  1,
			MaxAgeDays: 7,
			Compress:   true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()
	logger.Info().Msg("start cleanup")
	compressed := recent + ".gz"
	waitForRotationCleanup(t, "MaxAge removal and gzip compression", func() bool {
		_, staleErr := os.Stat(stale)
		_, plainErr := os.Stat(recent)
		info, archiveErr := os.Stat(compressed)
		return os.IsNotExist(staleErr) && os.IsNotExist(plainErr) && archiveErr == nil && info.Size() > 0
	})
	file, err := os.Open(compressed)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	contents, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "recent-data" {
		t.Fatalf("compressed backup content = %q, want recent-data", contents)
	}
}
