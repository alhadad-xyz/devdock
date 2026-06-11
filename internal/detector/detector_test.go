package detector_test

import (
	"path/filepath"
	"testing"

	"devdock/internal/detector"
)

func TestDetect(t *testing.T) {
	baseDir := "testdata/fixtures"

	tests := []struct {
		name     string
		fixture  string
		wantType string
		wantConf detector.Confidence
	}{
		// High confidence tests
		{"Laravel (High)", "laravel", "laravel", detector.High},
		{"Next.js JS (High)", "nextjs", "nextjs", detector.High},
		{"Next.js TS (High)", "nextjs-ts", "nextjs", detector.High},
		{"Docker Compose compose.yml (High)", "docker-compose", "docker-compose", detector.High},
		{"Docker Compose docker-compose.yml (High)", "docker-compose-yml", "docker-compose", detector.High},

		// Low confidence tests
		{"Laravel (Low)", "low-laravel", "laravel", detector.Low},
		{"Express Ambiguous (Low)", "low-nextjs", "express", detector.Low},

		// Unknown test
		{"Unknown Stack", "unknown", "unknown", detector.None},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectDir := filepath.Join(baseDir, tt.fixture)

			res := detector.Detect(projectDir)

			if res.Type != tt.wantType {
				t.Errorf("Detect() Type = %v, want %v", res.Type, tt.wantType)
			}
			if res.Confidence != tt.wantConf {
				t.Errorf("Detect() Confidence = %v, want %v", res.Confidence, tt.wantConf)
			}
		})
	}
}
