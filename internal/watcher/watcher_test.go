package watcher

import (
	"testing"

	"github.com/st0o0/vordr/internal/scanner"
)

func TestParseImageRef(t *testing.T) {
	tests := []struct {
		input string
		want  scanner.ImageRef
	}{
		{"grafana/grafana:13.3.0", scanner.ImageRef{Image: "grafana/grafana", Tag: "13.3.0"}},
		{"nginx", scanner.ImageRef{Image: "nginx", Tag: "latest"}},
		{"nginx:1.25", scanner.ImageRef{Image: "nginx", Tag: "1.25"}},
		{"ghcr.io/st0o0/vordr:1.0.0", scanner.ImageRef{Image: "ghcr.io/st0o0/vordr", Tag: "1.0.0"}},
		{"localhost:5000/myapp:dev", scanner.ImageRef{Image: "localhost:5000/myapp", Tag: "dev"}},
		{"registry.example.com:5000/org/app:v2", scanner.ImageRef{Image: "registry.example.com:5000/org/app", Tag: "v2"}},
		{"ubuntu", scanner.ImageRef{Image: "ubuntu", Tag: "latest"}},
		{"ghcr.io/org/app:rolling@sha256:52d09607f72d475ba75cfc781f9041cf679ddb16ed2e9dfb8a2139cc1e89ab94", scanner.ImageRef{Image: "ghcr.io/org/app", Tag: "rolling"}},
		{"nginx@sha256:abc123def456", scanner.ImageRef{Image: "nginx", Tag: "latest"}},
		{"ghcr.io/org/app:v2.1@sha256:deadbeef", scanner.ImageRef{Image: "ghcr.io/org/app", Tag: "v2.1"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseImageRef(tt.input)
			if got != tt.want {
				t.Errorf("ParseImageRef(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
