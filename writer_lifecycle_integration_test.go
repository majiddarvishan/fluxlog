package fluxlog

import (
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

type lifecycleBlockingWriter struct {
	entered chan struct{}
	release chan struct{}
}

func (w *lifecycleBlockingWriter) Write(p []byte) (int, error) {
	return w.WriteLevel(zerolog.NoLevel, p)
}

func (w *lifecycleBlockingWriter) WriteLevel(_ zerolog.Level, p []byte) (int, error) {
	select {
	case w.entered <- struct{}{}:
	default:
	}
	<-w.release
	return len(p), nil
}

type lifecycleCountingCloser struct {
	calls atomic.Int32
	err   error
}

func (c *lifecycleCountingCloser) Close() error {
	c.calls.Add(1)
	return c.err
}

func lifecycleReceive[T any](t *testing.T, ch <-chan T, label string) T {
	t.Helper()
	select {
	case result := <-ch:
		return result
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for %s", label)
		var zero T
		return zero
	}
}

// LIFE-01: Close must wait for an active write and close owned resources once.
func TestLIFE01CloseWaitsForActiveWrite(t *testing.T) {
	blocking := &lifecycleBlockingWriter{
		entered: make(chan struct{}, 1),
		release: make(chan struct{}),
	}
	closer := &lifecycleCountingCloser{}
	managed := &managedWriter{writer: blocking, closers: []io.Closer{closer}}
	type writeResult struct {
		count int
		err   error
	}
	wrote := make(chan writeResult, 1)
	go func() {
		n, err := managed.WriteLevel(zerolog.InfoLevel, []byte("event"))
		wrote <- writeResult{count: n, err: err}
	}()
	lifecycleReceive(t, blocking.entered, "writer entry")

	closed := make(chan error, 1)
	go func() { closed <- managed.Close() }()
	// The write still owns the lock. Close cannot have completed.
	select {
	case err := <-closed:
		t.Fatalf("Close finished before active write: %v", err)
	default:
	}
	close(blocking.release)
	result := lifecycleReceive(t, wrote, "active write")
	if result.count != len("event") || result.err != nil {
		t.Fatalf("write result = (%d, %v), want (%d, nil)", result.count, result.err, len("event"))
	}
	if err := lifecycleReceive(t, closed, "Close"); err != nil {
		t.Fatal(err)
	}
	if got := closer.calls.Load(); got != 1 {
		t.Fatalf("owned closer calls = %d, want 1", got)
	}
	if err := managed.Close(); err != nil {
		t.Fatal(err)
	}
	if got := closer.calls.Load(); got != 1 {
		t.Fatalf("owned closer calls after second Close = %d, want 1", got)
	}
	n, err := managed.WriteLevel(zerolog.InfoLevel, []byte("after"))
	if n != 0 || !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("post-close WriteLevel = (%d, %v), want (0, ErrClosedPipe)", n, err)
	}
}

// LIFE-01: concurrent writes may complete or report ErrClosedPipe at shutdown.
func TestLIFE01ConcurrentWritesAndClose(t *testing.T) {
	closer := &lifecycleCountingCloser{}
	managed := &managedWriter{
		writer:  zerolog.MultiLevelWriter(io.Discard),
		closers: []io.Closer{closer},
	}
	var wg sync.WaitGroup
	errorsSeen := make(chan error, 8)
	start := make(chan struct{})
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for i := 0; i < 150; i++ {
				n, err := managed.WriteLevel(zerolog.InfoLevel, []byte("data"))
				if err == nil && n == len("data") {
					continue
				}
				if errors.Is(err, io.ErrClosedPipe) && n == 0 {
					continue
				}
				errorsSeen <- errors.New("unexpected write count/error during Close")
				return
			}
		}()
	}
	close(start)
	if err := managed.Close(); err != nil {
		t.Fatal(err)
	}
	wg.Wait()
	close(errorsSeen)
	for err := range errorsSeen {
		t.Fatal(err)
	}
	if closer.calls.Load() != 1 {
		t.Fatalf("close count = %d, want 1", closer.calls.Load())
	}
}

// LIFE-01: the default logger pointer is independent of resource ownership.
func TestLIFE01SetDefaultDoesNotCloseReplacedLogger(t *testing.T) {
	first, err := New(Config{Level: InfoLevel, Console: &ConsoleConfig{Writer: io.Discard, Format: JSONFormat}})
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()
	second, err := New(Config{Level: InfoLevel, Console: &ConsoleConfig{Writer: io.Discard, Format: JSONFormat}})
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()

	previous := SetDefault(first)
	defer SetDefault(previous)
	if got := SetDefault(second); got != first {
		t.Fatalf("SetDefault returned previous = %p, want %p", got, first)
	}
	if _, err := first.output.WriteLevel(zerolog.InfoLevel, []byte("still-open")); err != nil {
		t.Fatalf("replaced logger was closed: %v", err)
	}

	var wg sync.WaitGroup
	for worker := 0; worker < 6; worker++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 150; i++ {
				if id%2 == 0 {
					SetDefault(first)
				} else {
					SetDefault(second)
				}
				Info("default swap event")
			}
		}(worker)
	}
	wg.Wait()
	if _, err := second.output.WriteLevel(zerolog.InfoLevel, []byte("still-open")); err != nil {
		t.Fatalf("replaced logger was closed by SetDefault: %v", err)
	}
}
