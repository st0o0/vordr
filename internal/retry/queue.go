package retry

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/st0o0/vordr/internal/scanner"
)

type Queue struct {
	scanCh     chan<- scanner.ImageRef
	maxBackoff time.Duration
	mu         sync.Mutex
	backoffs   map[string]time.Duration
	pending    chan retryItem
	logger     *slog.Logger
}

type retryItem struct {
	ref     scanner.ImageRef
	backoff time.Duration
}

const baseBackoff = 30 * time.Second

func New(scanCh chan<- scanner.ImageRef, maxBackoff time.Duration, logger *slog.Logger) *Queue {
	return &Queue{
		scanCh:     scanCh,
		maxBackoff: maxBackoff,
		backoffs:   make(map[string]time.Duration),
		pending:    make(chan retryItem, 100),
		logger:     logger,
	}
}

func (q *Queue) Add(ref scanner.ImageRef) {
	q.mu.Lock()
	key := ref.String()
	current := q.backoffs[key]
	if current == 0 {
		current = baseBackoff
	} else {
		current *= 2
		if current > q.maxBackoff {
			current = q.maxBackoff
		}
	}
	q.backoffs[key] = current
	q.mu.Unlock()

	q.logger.Info("scheduling retry", "event", "retry_scheduled", "image", key, "backoff", current)

	select {
	case q.pending <- retryItem{ref: ref, backoff: current}:
	default:
		q.logger.Warn("retry queue full, dropping", "event", "retry_dropped", "image", key)
	}
}

func (q *Queue) Reset(image string) {
	q.mu.Lock()
	delete(q.backoffs, image)
	q.mu.Unlock()
}

func (q *Queue) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case item := <-q.pending:
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(item.backoff):
				select {
				case <-ctx.Done():
					return nil
				case q.scanCh <- item.ref:
					q.logger.Debug("retry queued for scan", "event", "retry_sent", "image", item.ref.String())
				}
			}
		}
	}
}
