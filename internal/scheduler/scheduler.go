package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/st0o0/vordr/internal/scanner"
	"github.com/st0o0/vordr/internal/store"
)

type Scheduler struct {
	interval time.Duration
	store    *store.Store
	scanCh   chan<- scanner.ImageRef
	logger   *slog.Logger
}

func New(interval time.Duration, s *store.Store, scanCh chan<- scanner.ImageRef, logger *slog.Logger) *Scheduler {
	return &Scheduler{
		interval: interval,
		store:    s,
		scanCh:   scanCh,
		logger:   logger,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			s.rescan(ctx)
		}
	}
}

func (s *Scheduler) rescan(ctx context.Context) {
	images := s.store.Images()
	s.logger.Info("scheduled rescan", "event", "scheduled_rescan", "images", len(images))

	for _, ref := range images {
		select {
		case <-ctx.Done():
			return
		case s.scanCh <- ref:
		}
	}
}
