package ui_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/st0o0/vordr/internal/ui"
)

func TestHandler_ServesHTML(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ui.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("expected Content-Type text/html, got %q", ct)
	}

	if rec.Body.Len() == 0 {
		t.Fatal("expected non-empty body")
	}

	if !strings.Contains(rec.Body.String(), "<!DOCTYPE html>") {
		t.Fatal("expected HTML document in body")
	}
}
