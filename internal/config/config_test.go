package config

import (
	"reflect"
	"testing"
)

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestShouldExclude(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		image    string
		want     bool
	}{
		{"exact match", []string{"nginx:latest"}, "nginx:latest", true},
		{"non-matching", []string{"nginx:latest"}, "redis:latest", false},
		{"empty patterns", nil, "nginx:latest", false},
		{"multiple patterns", []string{"redis:*", "nginx:latest"}, "nginx:latest", true},
		{"wildcard tag", []string{"nginx:*"}, "nginx:1.25", true},
		{"wildcard no match", []string{"nginx:*"}, "redis:7", false},
		{"wildcard name", []string{"*:latest"}, "app:latest", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{ExcludeImages: tt.patterns}
			if got := c.ShouldExclude(tt.image); got != tt.want {
				t.Errorf("ShouldExclude(%q) = %v, want %v", tt.image, got, tt.want)
			}
		})
	}
}

func TestLoadExcludeImages(t *testing.T) {
	cfg, err := Load(env(map[string]string{"VORDR_EXCLUDE_IMAGES": "nginx:*, redis:7 ,, postgres:*"}))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"nginx:*", "redis:7", "postgres:*"}
	if !reflect.DeepEqual(cfg.ExcludeImages, want) {
		t.Errorf("got %v, want %v", cfg.ExcludeImages, want)
	}

	cfg, err = Load(env(nil))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.ExcludeImages) != 0 {
		t.Errorf("expected empty, got %v", cfg.ExcludeImages)
	}
}

func TestLoadTrivyServerURL(t *testing.T) {
	cfg, err := Load(env(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TrivyServerURL != "" {
		t.Errorf("default = %q, want empty", cfg.TrivyServerURL)
	}

	cfg, err = Load(env(map[string]string{"VORDR_TRIVY_SERVER_URL": "http://trivy:4954"}))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TrivyServerURL != "http://trivy:4954" {
		t.Errorf("got %q, want http://trivy:4954", cfg.TrivyServerURL)
	}
}

func TestLoadAlertSeverity(t *testing.T) {
	cfg, err := Load(env(nil))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AlertSeverity != "HIGH" {
		t.Errorf("default = %q, want HIGH", cfg.AlertSeverity)
	}

	for _, v := range []string{"unknown", "LOW", "Medium", "high", "CRITICAL"} {
		cfg, err := Load(env(map[string]string{"VORDR_ALERT_SEVERITY": v}))
		if err != nil {
			t.Errorf("%q: unexpected error %v", v, err)
			continue
		}
		if cfg.AlertSeverity == "" || cfg.AlertSeverity != toUpper(v) {
			t.Errorf("%q: got %q", v, cfg.AlertSeverity)
		}
	}

	if _, err := Load(env(map[string]string{"VORDR_ALERT_SEVERITY": "bogus"})); err == nil {
		t.Error("expected error for invalid severity")
	}
}

func toUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}
