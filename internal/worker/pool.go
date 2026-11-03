package worker

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/st0o0/vordr/internal/scanner"
	"github.com/st0o0/vordr/internal/store"
	"github.com/st0o0/vordr/internal/trivy"
)

type Scanner interface {
	Scan(ctx context.Context, imageRef string) (*trivy.ScanResult, error)
}

type Pool struct {
	workers       int
	scanCh        <-chan scanner.ImageRef
	store         *store.Store
	scanner       Scanner
	ActiveWorkers atomic.Int64
	logger        *slog.Logger
}

func NewPool(workers int, scanCh <-chan scanner.ImageRef, s *store.Store, sc Scanner, logger *slog.Logger) *Pool {
	return &Pool{
		workers: workers,
		scanCh:  scanCh,
		store:   s,
		scanner: sc,
		logger:  logger,
	}
}

func (p *Pool) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for range p.workers {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return nil
				case ref, ok := <-p.scanCh:
					if !ok {
						return nil
					}
					p.scan(ctx, ref)
				}
			}
		})
	}

	return g.Wait()
}

func (p *Pool) scan(ctx context.Context, ref scanner.ImageRef) {
	key := ref.String()

	if !p.store.SetScanning(key) {
		p.logger.Debug("skipping already-scanning image", "event", "scan_skipped", "image", key)
		return
	}

	p.ActiveWorkers.Add(1)
	defer p.ActiveWorkers.Add(-1)

	p.logger.Info("scanning image", "event", "scan_started", "image", key)
	start := time.Now()

	result, err := p.scanner.Scan(ctx, ref.ScanTarget())
	duration := time.Since(start)

	if err != nil {
		if ctx.Err() != nil {
			return
		}
		p.logger.Error("scan failed", "event", "scan_failed", "image", key, "error", err, "duration", duration)
		p.store.SetError(key, err.Error())
		return
	}

	vulns := make([]scanner.Vulnerability, len(result.Vulnerabilities))
	for i, v := range result.Vulnerabilities {
		vulns[i] = scanner.Vulnerability(v)
	}

	p.store.UpdateResult(key, result.Counts, result.Fixable, vulns, duration)
	p.logger.Info("scan complete", "event", "scan_completed", "image", key, "duration", duration, "counts", result.Counts)
}
