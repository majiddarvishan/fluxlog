package fluxlog

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

type outputFailingWriter struct {
	err error
}

func (w outputFailingWriter) Write(_ []byte) (int, error) {
	return 0, w.err
}

type outputShortWriter struct{}

func (outputShortWriter) Write(p []byte) (int, error) {
	return len(p) - 1, nil
}

// ERR-01: writer errors are observable at the managed output boundary.
func TestERR01CustomWriterFailureAndShortWrite(t *testing.T) {
	sentinel := errors.New("sink failure")
	for _, tc := range []struct {
		name string
		writer io.Writer
		want error
	}{
		{name: "error", writer: outputFailingWriter{err: sentinel}, want: sentinel},
		{name: "short-write", writer: outputShortWriter{}, want: io.ErrShortWrite},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger, err := New(Config{
				Level: InfoLevel,
				Console: &ConsoleConfig{Writer: tc.writer, Format: JSONFormat},
			})
			if err != nil {
				t.Fatal(err)
			}
			defer logger.Close()
			n, err := logger.output.WriteLevel(zerolog.InfoLevel, []byte("record"))
			if !errors.Is(err, tc.want) {
				t.Fatalf("write error = %v, want %v (n=%d)", err, tc.want, n)
			}
		})
	}
}

// ERR-01: a broken stream must not prevent a later file sink from receiving
// the same event; however the first sink error is returned to the caller.
func TestERR01MultiSinkPartialFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "application.log")
	sentinel := errors.New("console down")
	logger, err := New(Config{
		Level: InfoLevel,
		Console: &ConsoleConfig{Writer: outputFailingWriter{err: sentinel}, Format: JSONFormat},
		File: &FileConfig{Path: path, Format: JSONFormat},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer logger.Close()
	record := []byte("{\"level\":\"info\",\"message\":\"record\"}\n")
	_, err = logger.output.WriteLevel(zerolog.InfoLevel, record)
	if !errors.Is(err, sentinel) {
		t.Fatalf("first sink error was lost: %v", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(content, record) {
		t.Fatalf("file sink content = %q, want %q", content, record)
	}
}

// ERR-01: use a regular file as a directory to fail reliably without chmod,
// root permissions, or assumptions about a platform's permission policy.
func TestERR01CreateLogDirectoryFailure(t *testing.T) {
	base := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(base, []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}
	logger, err := New(Config{Level: InfoLevel, File: &FileConfig{
		Path: filepath.Join(base, "log.json"),
		Format: JSONFormat,
	}})
	if err == nil {
		if logger != nil {
			logger.Close()
		}
		t.Fatal("New unexpectedly succeeded with a regular file as parent directory")
	}
	if !strings.Contains(err.Error(), "create log directory") {
		t.Fatalf("unexpected directory error: %v", err)
	}
}
