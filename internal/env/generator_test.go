package env_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"devdock/internal/config"
	"devdock/internal/env"
)

func TestEnvGenerateLaravel(t *testing.T) {
	tmpDir := t.TempDir()
	
	cfg := &config.Config{
		Project: config.ProjectConfig{Name: "my app", Type: "laravel"},
	}
	cfg.Services.MySQL = &config.MySQLConfig{Enabled: true}
	
	err := env.Generate(tmpDir, cfg)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	envPath := filepath.Join(tmpDir, ".env")
	b, _ := os.ReadFile(envPath)
	content := string(b)
	
	if !strings.Contains(content, "DB_DATABASE=my_app") {
		t.Errorf("Missing DB_DATABASE in .env")
	}
}

func TestEnvGenerateNextJS(t *testing.T) {
	tmpDir := t.TempDir()
	
	cfg := &config.Config{
		Project: config.ProjectConfig{Name: "my app", Type: "nextjs"},
	}
	cfg.Services.Postgres = &config.PostgresConfig{Enabled: true}
	
	err := env.Generate(tmpDir, cfg)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	envPath := filepath.Join(tmpDir, ".env")
	b, _ := os.ReadFile(envPath)
	content := string(b)
	
	if !strings.Contains(content, "DATABASE_URL=postgresql://devdock:password@postgres:5432/my_app") {
		t.Errorf("Missing DATABASE_URL in .env")
	}
}
