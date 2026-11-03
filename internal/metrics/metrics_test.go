package metrics

import (
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/st0o0/vordr/internal/scanner"
)

type mockStore struct {
	states []scanner.ImageState
}

func (m *mockStore) Snapshot() []scanner.ImageState {
	return m.states
}

func TestCollectorEmptyStore(t *testing.T) {
	var active atomic.Int64
	c := NewCollector(&mockStore{}, &active, func() int { return 0 })

	registry := prometheus.NewRegistry()
	registry.MustRegister(c)

	count, err := testutil.GatherAndCount(registry)
	if err != nil {
		t.Fatal(err)
	}
	// Should have 6 operational metrics even with empty store
	if count != 6 {
		t.Fatalf("expected 6 metrics (operational), got %d", count)
	}
}

func TestCollectorWithScannedImage(t *testing.T) {
	store := &mockStore{
		states: []scanner.ImageState{
			{
				Ref:        scanner.ImageRef{Image: "nginx", Tag: "latest"},
				Status:     scanner.Scanned,
				Containers: []string{"web"},
				Counts:     map[string]int{"HIGH": 3, "MEDIUM": 12},
				Fixable:    map[string]int{"HIGH": 2},
				ScannedAt:  time.Unix(1730000000, 0),
				Duration:   14200 * time.Millisecond,
			},
		},
	}

	var active atomic.Int64
	active.Store(1)
	c := NewCollector(store, &active, func() int { return 2 })

	registry := prometheus.NewRegistry()
	registry.MustRegister(c)

	expected := `
# HELP vordr_scan_workers_active Currently active scan workers
# TYPE vordr_scan_workers_active gauge
vordr_scan_workers_active 1
`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "vordr_scan_workers_active"); err != nil {
		t.Fatal(err)
	}

	expected = `
# HELP vordr_images_scanned Images successfully scanned
# TYPE vordr_images_scanned gauge
vordr_images_scanned 1
`
	if err := testutil.CollectAndCompare(c, strings.NewReader(expected), "vordr_images_scanned"); err != nil {
		t.Fatal(err)
	}
}

func TestCollectorPendingImageNoVulnMetrics(t *testing.T) {
	store := &mockStore{
		states: []scanner.ImageState{
			{
				Ref:        scanner.ImageRef{Image: "nginx", Tag: "latest"},
				Status:     scanner.Pending,
				Containers: []string{"web"},
			},
		},
	}

	var active atomic.Int64
	c := NewCollector(store, &active, func() int { return 0 })

	registry := prometheus.NewRegistry()
	registry.MustRegister(c)

	metrics, err := testutil.GatherAndCount(registry, "vordr_vulnerability_count")
	if err != nil {
		t.Fatal(err)
	}
	if metrics != 0 {
		t.Fatalf("expected 0 vulnerability metrics for pending image, got %d", metrics)
	}
}
