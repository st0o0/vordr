package store

import (
	"slices"
	"sync"
	"time"

	"github.com/st0o0/vordr/internal/scanner"
)

type Store struct {
	mu     sync.RWMutex
	images map[string]*entry
}

type entry struct {
	ref        scanner.ImageRef
	status     scanner.Status
	containers []string
	counts     map[string]int
	fixable    map[string]int
	scannedAt  time.Time
	duration   time.Duration
	errMsg     string
}

func New() *Store {
	return &Store{images: make(map[string]*entry)}
}

func (s *Store) Add(ref scanner.ImageRef, container string) (isNew bool) {
	key := ref.String()
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.images[key]
	if !exists {
		s.images[key] = &entry{
			ref:        ref,
			status:     scanner.Pending,
			containers: []string{container},
		}
		return true
	}

	if !slices.Contains(e.containers, container) {
		e.containers = append(e.containers, container)
	}
	return false
}

func (s *Store) Remove(container string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, e := range s.images {
		idx := slices.Index(e.containers, container)
		if idx < 0 {
			continue
		}
		e.containers = slices.Delete(e.containers, idx, idx+1)
		if len(e.containers) == 0 {
			delete(s.images, key)
		}
		return
	}
}

func (s *Store) SetScanning(image string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.images[image]
	if !ok {
		return false
	}
	if e.status == scanner.Scanning {
		return false
	}
	e.status = scanner.Scanning
	return true
}

func (s *Store) UpdateResult(image string, counts, fixable map[string]int, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.images[image]
	if !ok {
		return
	}
	e.status = scanner.Scanned
	e.counts = counts
	e.fixable = fixable
	e.scannedAt = time.Now()
	e.duration = duration
	e.errMsg = ""
}

func (s *Store) SetError(image string, errMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.images[image]
	if !ok {
		return
	}
	e.status = scanner.Error
	e.errMsg = errMsg
}

func (s *Store) Snapshot() []scanner.ImageState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]scanner.ImageState, 0, len(s.images))
	for _, e := range s.images {
		containers := make([]string, len(e.containers))
		copy(containers, e.containers)

		state := scanner.ImageState{
			Ref:        e.ref,
			Status:     e.status,
			Containers: containers,
			ScannedAt:  e.scannedAt,
			Duration:   e.duration,
			Error:      e.errMsg,
		}

		if e.counts != nil {
			state.Counts = make(map[string]int, len(e.counts))
			for k, v := range e.counts {
				state.Counts[k] = v
			}
		}
		if e.fixable != nil {
			state.Fixable = make(map[string]int, len(e.fixable))
			for k, v := range e.fixable {
				state.Fixable[k] = v
			}
		}

		out = append(out, state)
	}
	return out
}

func (s *Store) Images() []scanner.ImageRef {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]scanner.ImageRef, 0, len(s.images))
	for _, e := range s.images {
		out = append(out, e.ref)
	}
	return out
}
