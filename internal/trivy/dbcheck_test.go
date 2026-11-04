package trivy

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"
)

func TestDBCheckerDisabledReturnsImmediately(t *testing.T) {
	checker := NewDBChecker("/tmp", time.Hour, slog.Default(), false)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := checker.Run(ctx)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if checker.DBUpdatedAt() != (time.Time{}) {
		t.Fatal("expected zero time when disabled")
	}
}

func TestParseVersionOutput(t *testing.T) {
	raw := `{"Version":"0.76.1","VulnerabilityDB":{"Version":2,"UpdatedAt":"2026-11-03T13:13:42.304243913Z"}}`
	var ver versionOutput
	if err := json.Unmarshal([]byte(raw), &ver); err != nil {
		t.Fatal(err)
	}
	if ver.VulnerabilityDB.UpdatedAt.IsZero() {
		t.Fatal("expected non-zero UpdatedAt")
	}
}
