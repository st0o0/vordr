package trivy

import (
	"encoding/json"
	"testing"
)

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
