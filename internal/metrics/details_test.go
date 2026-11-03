package metrics

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/st0o0/vordr/internal/scanner"
)

func TestDetailsHandlerAllImages(t *testing.T) {
	store := &mockStore{
		states: []scanner.ImageState{
			{
				Ref:    scanner.ImageRef{Image: "nginx", Tag: "latest"},
				Status: scanner.Scanned,
				Vulnerabilities: []scanner.Vulnerability{
					{ID: "CVE-2024-0001", PkgName: "libssl", InstalledVersion: "1.1.1", FixedVersion: "1.1.2", Severity: "HIGH"},
				},
				ScannedAt: time.Unix(1730000000, 0),
			},
			{
				Ref:    scanner.ImageRef{Image: "alpine", Tag: "3.19"},
				Status: scanner.Pending,
			},
		},
	}

	handler := detailsHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/details", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp []detailsResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if len(resp) != 2 {
		t.Fatalf("expected 2 images, got %d", len(resp))
	}
}

func TestDetailsHandlerFilterByImage(t *testing.T) {
	store := &mockStore{
		states: []scanner.ImageState{
			{
				Ref:    scanner.ImageRef{Image: "nginx", Tag: "latest"},
				Status: scanner.Scanned,
			},
			{
				Ref:    scanner.ImageRef{Image: "alpine", Tag: "3.19"},
				Status: scanner.Pending,
			},
		},
	}

	handler := detailsHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/details?image=nginx:latest", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp []detailsResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if len(resp) != 1 {
		t.Fatalf("expected 1 image, got %d", len(resp))
	}
	if resp[0].Image != "nginx" {
		t.Fatalf("expected nginx, got %s", resp[0].Image)
	}
}

func TestDetailsHandlerNotFound(t *testing.T) {
	store := &mockStore{states: []scanner.ImageState{}}

	handler := detailsHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/details?image=nonexistent:latest", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestDetailsHandlerEmptyStore(t *testing.T) {
	store := &mockStore{}

	handler := detailsHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/details", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
