package worker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/st0o0/vordr/internal/scanner"
	"github.com/st0o0/vordr/internal/store"
	"github.com/st0o0/vordr/internal/trivy"
)

type mockScanner struct {
	result *trivy.ScanResult
	err    error
}

func (m *mockScanner) Scan(_ context.Context, _ string) (*trivy.ScanResult, error) {
	return m.result, m.err
}

func TestPoolScanSuccess(t *testing.T) {
	s := store.New()
	ref := scanner.ImageRef{Image: "nginx", Tag: "latest"}
	s.Add(ref, "web")

	scanCh := make(chan scanner.ImageRef, 1)
	scanCh <- ref

	mock := &mockScanner{result: &trivy.ScanResult{
		Counts:  map[string]int{"HIGH": 1},
		Fixable: map[string]int{"HIGH": 1},
	}}

	pool := NewPool(1, scanCh, s, mock, discardLogger())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		time.Sleep(500 * time.Millisecond)
		close(scanCh)
	}()

	if err := pool.Run(ctx); err != nil {
		t.Fatal(err)
	}

	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 image, got %d", len(snap))
	}
	if snap[0].Status != scanner.Scanned {
		t.Fatalf("expected Scanned, got %v", snap[0].Status)
	}
	if snap[0].Counts["HIGH"] != 1 {
		t.Fatalf("expected HIGH=1, got %d", snap[0].Counts["HIGH"])
	}
}

func TestPoolScanError(t *testing.T) {
	s := store.New()
	ref := scanner.ImageRef{Image: "nginx", Tag: "latest"}
	s.Add(ref, "web")

	scanCh := make(chan scanner.ImageRef, 1)
	scanCh <- ref

	mock := &mockScanner{err: fmt.Errorf("trivy failed")}

	pool := NewPool(1, scanCh, s, mock, discardLogger())

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		time.Sleep(500 * time.Millisecond)
		close(scanCh)
	}()

	if err := pool.Run(ctx); err != nil {
		t.Fatal(err)
	}

	snap := s.Snapshot()
	if snap[0].Status != scanner.Error {
		t.Fatalf("expected Error, got %v", snap[0].Status)
	}
}
