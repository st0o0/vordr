package store

import (
	"sync"
	"testing"
	"time"

	"github.com/st0o0/vordr/internal/scanner"
)

func ref(image, tag string) scanner.ImageRef {
	return scanner.ImageRef{Image: image, Tag: tag}
}

func TestAddNewImage(t *testing.T) {
	s := New()
	isNew := s.Add(ref("nginx", "latest"), "web")
	if !isNew {
		t.Fatal("expected isNew=true for first add")
	}

	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 image, got %d", len(snap))
	}
	if snap[0].Status != scanner.Pending {
		t.Fatalf("expected Pending, got %v", snap[0].Status)
	}
	if len(snap[0].Containers) != 1 || snap[0].Containers[0] != "web" {
		t.Fatalf("expected [web], got %v", snap[0].Containers)
	}
}

func TestAddExistingImage(t *testing.T) {
	s := New()
	s.Add(ref("nginx", "latest"), "web1")
	isNew := s.Add(ref("nginx", "latest"), "web2")
	if isNew {
		t.Fatal("expected isNew=false for duplicate image")
	}

	snap := s.Snapshot()
	if len(snap[0].Containers) != 2 {
		t.Fatalf("expected 2 containers, got %d", len(snap[0].Containers))
	}
}

func TestAddDuplicateContainer(t *testing.T) {
	s := New()
	s.Add(ref("nginx", "latest"), "web")
	s.Add(ref("nginx", "latest"), "web")

	snap := s.Snapshot()
	if len(snap[0].Containers) != 1 {
		t.Fatalf("expected 1 container (deduped), got %d", len(snap[0].Containers))
	}
}

func TestRemoveLastContainer(t *testing.T) {
	s := New()
	s.Add(ref("nginx", "latest"), "web")
	s.Remove("web")

	snap := s.Snapshot()
	if len(snap) != 0 {
		t.Fatalf("expected 0 images after removing last container, got %d", len(snap))
	}
}

func TestRemoveOneOfMultiple(t *testing.T) {
	s := New()
	s.Add(ref("nginx", "latest"), "web1")
	s.Add(ref("nginx", "latest"), "web2")
	s.Remove("web1")

	snap := s.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 image, got %d", len(snap))
	}
	if len(snap[0].Containers) != 1 || snap[0].Containers[0] != "web2" {
		t.Fatalf("expected [web2], got %v", snap[0].Containers)
	}
}

func TestSetScanningAndUpdate(t *testing.T) {
	s := New()
	s.Add(ref("nginx", "latest"), "web")

	ok := s.SetScanning("nginx:latest")
	if !ok {
		t.Fatal("expected SetScanning to return true")
	}

	ok = s.SetScanning("nginx:latest")
	if ok {
		t.Fatal("expected SetScanning to return false for already-scanning image")
	}

	counts := map[string]int{"HIGH": 3, "MEDIUM": 12}
	fixable := map[string]int{"HIGH": 2}
	vulns := []scanner.Vulnerability{
		{ID: "CVE-2024-0001", PkgName: "libfoo", InstalledVersion: "1.0", FixedVersion: "1.1", Severity: "HIGH"},
	}
	s.UpdateResult("nginx:latest", counts, fixable, vulns, 14*time.Second)

	snap := s.Snapshot()
	if len(snap[0].Vulnerabilities) != 1 || snap[0].Vulnerabilities[0] != vulns[0] {
		t.Fatalf("unexpected vulnerabilities: %+v", snap[0].Vulnerabilities)
	}
	if snap[0].Status != scanner.Scanned {
		t.Fatalf("expected Scanned, got %v", snap[0].Status)
	}
	if snap[0].Counts["HIGH"] != 3 {
		t.Fatalf("expected HIGH=3, got %d", snap[0].Counts["HIGH"])
	}
	if snap[0].Fixable["HIGH"] != 2 {
		t.Fatalf("expected fixable HIGH=2, got %d", snap[0].Fixable["HIGH"])
	}
	if snap[0].PreviousVulnIDs != nil {
		t.Fatalf("expected nil PreviousVulnIDs after first scan, got %v", snap[0].PreviousVulnIDs)
	}

	vulns2 := []scanner.Vulnerability{
		{ID: "CVE-2024-0002", Severity: "CRITICAL"},
	}
	s.UpdateResult("nginx:latest", counts, fixable, vulns2, time.Second)
	snap = s.Snapshot()
	if len(snap[0].PreviousVulnIDs) != 1 || !snap[0].PreviousVulnIDs["CVE-2024-0001"] {
		t.Fatalf("unexpected PreviousVulnIDs: %v", snap[0].PreviousVulnIDs)
	}
}

func TestSetError(t *testing.T) {
	s := New()
	s.Add(ref("nginx", "latest"), "web")
	s.SetScanning("nginx:latest")
	s.SetError("nginx:latest", "trivy unreachable")

	snap := s.Snapshot()
	if snap[0].Status != scanner.Error {
		t.Fatalf("expected Error, got %v", snap[0].Status)
	}
	if snap[0].Error != "trivy unreachable" {
		t.Fatalf("expected error message, got %q", snap[0].Error)
	}
}

func TestSnapshotIsolation(t *testing.T) {
	s := New()
	s.Add(ref("nginx", "latest"), "web")

	snap := s.Snapshot()
	snap[0].Containers = append(snap[0].Containers, "hacked")

	snap2 := s.Snapshot()
	if len(snap2[0].Containers) != 1 {
		t.Fatal("snapshot mutation leaked into store")
	}
}

func TestImages(t *testing.T) {
	s := New()
	s.Add(ref("nginx", "latest"), "web")
	s.Add(ref("grafana/grafana", "13.3.0"), "grafana")

	images := s.Images()
	if len(images) != 2 {
		t.Fatalf("expected 2 images, got %d", len(images))
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := New()
	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r := ref("nginx", "latest")
			s.Add(r, "web")
			s.Snapshot()
			if i%2 == 0 {
				s.SetScanning("nginx:latest")
			}
		}(i)
	}

	wg.Wait()
}
