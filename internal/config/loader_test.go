package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"devdock/internal/config"
)

func TestLoader(t *testing.T) {
	tempDir := t.TempDir()

	writeConfig := func(content string) {
		_ = os.WriteFile(filepath.Join(tempDir, ".devdock.yml"), []byte(content), 0644)
	}

	// Valid Next.js
	writeConfig(`
version: "1"
project:
  name: my-next
  type: nextjs
app:
  command: "pnpm dev"
  port: 3000
`)
	if _, err := config.Load(tempDir); err != nil {
		t.Fatalf("expected valid nextjs config, got: %v", err)
	}

	// Valid Laravel
	writeConfig(`
version: "1"
project:
  name: my-laravel
  type: laravel
app:
  command: "php artisan serve"
  port: 8000
`)
	if _, err := config.Load(tempDir); err != nil {
		t.Fatalf("expected valid laravel config, got: %v", err)
	}

	// Valid Docker Compose
	writeConfig(`
version: "1"
project:
  name: my-dc
  type: docker-compose
`)
	if _, err := config.Load(tempDir); err != nil {
		t.Fatalf("expected valid docker-compose config, got: %v", err)
	}

	// Invalid Version
	writeConfig(`
version: "2"
project:
  name: test
  type: nextjs
app:
  command: "dev"
  port: 3000
`)
	if _, err := config.Load(tempDir); err == nil {
		t.Fatal("expected error for invalid version")
	}

	// Invalid Project Name
	writeConfig(`
version: "1"
project:
  name: Invalid Name!
  type: nextjs
app:
  command: "dev"
  port: 3000
`)
	if _, err := config.Load(tempDir); err == nil {
		t.Fatal("expected error for invalid project name")
	}

	// Invalid Project Type
	writeConfig(`
version: "1"
project:
  name: test
  type: unknown
app:
  command: "dev"
  port: 3000
`)
	if _, err := config.Load(tempDir); err == nil {
		t.Fatal("expected error for invalid project type")
	}

	// Missing app.command on nextjs
	writeConfig(`
version: "1"
project:
  name: test
  type: nextjs
app:
  port: 3000
`)
	if _, err := config.Load(tempDir); err == nil {
		t.Fatal("expected error for missing command on nextjs")
	}

	// Missing app.port on laravel
	writeConfig(`
version: "1"
project:
  name: test
  type: laravel
app:
  command: "serve"
`)
	if _, err := config.Load(tempDir); err == nil {
		t.Fatal("expected error for missing port on laravel")
	}
}

func TestBackwardCompatibility(t *testing.T) {
	fixtures := []string{
		"v0.1-laravel.yml",
		"v0.1-nextjs.yml",
		"v0.1-docker-compose.yml",
	}

	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			b, err := os.ReadFile(filepath.Join("testdata", "backward-compat", fixture))
			if err != nil {
				t.Fatalf("failed to read fixture: %v", err)
			}
			
			tempDir := t.TempDir()
			err = os.WriteFile(filepath.Join(tempDir, ".devdock.yml"), b, 0644)
			if err != nil {
				t.Fatalf("failed to write fixture: %v", err)
			}
			
			cfg, err := config.Load(tempDir)
			if err != nil {
				t.Fatalf("expected valid config for %s, got error: %v", fixture, err)
			}
			
			if cfg.Version != "1" {
				t.Errorf("expected version 1, got %s", cfg.Version)
			}
		})
	}
}
