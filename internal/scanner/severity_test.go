package scanner

import "testing"

func TestMeetsThreshold(t *testing.T) {
	tests := []struct {
		severity, threshold string
		want                bool
	}{
		{"CRITICAL", "HIGH", true},
		{"HIGH", "HIGH", true},
		{"MEDIUM", "HIGH", false},
		{"LOW", "HIGH", false},
		{"UNKNOWN", "LOW", false},
		{"CRITICAL", "CRITICAL", true},
		{"LOW", "LOW", true},
	}
	for _, tt := range tests {
		t.Run(tt.severity+"_"+tt.threshold, func(t *testing.T) {
			if got := MeetsThreshold(tt.severity, tt.threshold); got != tt.want {
				t.Errorf("MeetsThreshold(%q, %q) = %v, want %v", tt.severity, tt.threshold, got, tt.want)
			}
		})
	}
}
