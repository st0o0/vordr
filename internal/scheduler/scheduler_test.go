package scheduler

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/st0o0/vordr/internal/scanner"
	"github.com/st0o0/vordr/internal/store"
)

func TestSchedulerRescan(t *testing.T) {
	s := store.New()
	s.Add(scanner.ImageRef{Image: "nginx", Tag: "latest"}, "web")
	s.Add(scanner.ImageRef{Image: "grafana/grafana", Tag: "13.3.0"}, "grafana")

	scanCh := make(chan scanner.ImageRef, 10)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	sched := New(100*time.Millisecond, s, scanCh, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	_ = sched.Run(ctx)

	if len(scanCh) < 2 {
		t.Fatalf("expected at least 2 images queued for rescan, got %d", len(scanCh))
	}
}
