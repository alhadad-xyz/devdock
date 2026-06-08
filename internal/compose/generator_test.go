package compose_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"devdock/internal/compose"
	"devdock/internal/config"
)

func TestComposeGenerate(t *testing.T) {
	tmpDir := t.TempDir()
	
	cfg := &config.Config{
		Project: config.ProjectConfig{Name: "my app"},
	}
	cfg.Services.Postgres = &config.PostgresConfig{
		Enabled: true,
		Version: "15",
		Port:    5432,
	}
	cfg.Services.Redis = &config.RedisConfig{
		Enabled: true,
		Version: "7",
		Port:    6379,
	}
	
	err := compose.Generate(tmpDir, cfg)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	composePath := filepath.Join(tmpDir, "compose.yml")
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		t.Fatalf("compose.yml not created")
	}

	b, _ := os.ReadFile(composePath)
	content := string(b)
	
	if !compose.IsDevDockOwned(composePath) {
		t.Errorf("IsDevDockOwned failed, expected true. Content: %s", content)
	}

	if !strings.Contains(content, "postgres:15") {
		t.Errorf("Missing postgres image in compose.yml")
	}
	if !strings.Contains(content, "redis:7") {
		t.Errorf("Missing redis image in compose.yml")
	}
	if !strings.Contains(content, "devdock_my-app_postgres_data:/var/lib/postgresql/data") {
		t.Errorf("Missing postgres volume mapping in compose.yml")
	}
}
