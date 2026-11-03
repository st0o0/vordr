package retry

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/st0o0/vordr/internal/scanner"
)

func TestBackoffEscalation(t *testing.T) {
	scanCh := make(chan scanner.ImageRef, 10)
	q := New(scanCh, 15*time.Minute, discardLogger())
	ref := scanner.ImageRef{Image: "nginx", Tag: "latest"}

	q.Add(ref)
	q.mu.Lock()
	if q.backoffs["nginx:latest"] != 30*time.Second {
		t.Fatalf("expected 30s, got %v", q.backoffs["nginx:latest"])
	}
	q.mu.Unlock()

	q.Add(ref)
	q.mu.Lock()
	if q.backoffs["nginx:latest"] != 60*time.Second {
		t.Fatalf("expected 60s, got %v", q.backoffs["nginx:latest"])
	}
	q.mu.Unlock()

	q.Reset("nginx:latest")
	q.mu.Lock()
	if _, ok := q.backoffs["nginx:latest"]; ok {
		t.Fatal("expected backoff to be cleared after reset")
	}
	q.mu.Unlock()
}

func TestBackoffMaxCap(t *testing.T) {
	scanCh := make(chan scanner.ImageRef, 10)
	q := New(scanCh, 2*time.Minute, discardLogger())
	ref := scanner.ImageRef{Image: "nginx", Tag: "latest"}

	for range 10 {
		q.Add(ref)
	}
	q.mu.Lock()
	if q.backoffs["nginx:latest"] > 2*time.Minute {
		t.Fatalf("exceeded max backoff: %v", q.backoffs["nginx:latest"])
	}
	q.mu.Unlock()
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
