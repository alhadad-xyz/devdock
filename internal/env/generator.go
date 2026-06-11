package env

import (
	"bytes"
	"devdock/internal/config"
	"devdock/internal/utils"
	"fmt"
	"os"
	"path/filepath"
)

func Generate(projectDir string, cfg *config.Config) error {
	var buf bytes.Buffer

	dbName := utils.NormalizeDBName(cfg.Project.Name)

	if cfg.Project.Type == "laravel" {
		if cfg.Services.MySQL != nil && cfg.Services.MySQL.Enabled {
			buf.WriteString(fmt.Sprintf("DB_CONNECTION=mysql\n"))
			buf.WriteString(fmt.Sprintf("DB_HOST=mysql\n"))
			buf.WriteString(fmt.Sprintf("DB_PORT=3306\n"))
			buf.WriteString(fmt.Sprintf("DB_DATABASE=%s\n", dbName))
			buf.WriteString(fmt.Sprintf("DB_USERNAME=devdock\n"))
			buf.WriteString(fmt.Sprintf("DB_PASSWORD=password\n"))
		}
		if cfg.Services.Redis != nil && cfg.Services.Redis.Enabled {
			buf.WriteString(fmt.Sprintf("REDIS_HOST=redis\n"))
			buf.WriteString(fmt.Sprintf("REDIS_PASSWORD=null\n"))
			buf.WriteString(fmt.Sprintf("REDIS_PORT=6379\n"))
		}
	} else if cfg.Project.Type == "nextjs" || cfg.Project.Type == "express" || cfg.Project.Type == "fiber" {
		if cfg.Services.Postgres != nil && cfg.Services.Postgres.Enabled {
			buf.WriteString(fmt.Sprintf("DATABASE_URL=postgresql://devdock:password@postgres:5432/%s\n", dbName))
		}
		if cfg.Services.Redis != nil && cfg.Services.Redis.Enabled {
			buf.WriteString(fmt.Sprintf("REDIS_URL=redis://redis:6379\n"))
		}
	}

	envStr := buf.String()
	if envStr == "" {
		return nil // nothing to write
	}

	envPath := filepath.Join(projectDir, ".env")
	envExamplePath := filepath.Join(projectDir, ".env.example")

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		os.WriteFile(envPath, []byte(envStr), 0644)
	}

	if _, err := os.Stat(envExamplePath); os.IsNotExist(err) {
		os.WriteFile(envExamplePath, []byte(envStr), 0644)
	}

	updateGitignore(projectDir)
	return nil
}

func updateGitignore(projectDir string) {
	gitignorePath := filepath.Join(projectDir, ".gitignore")
	b, err := os.ReadFile(gitignorePath)
	if err != nil {
		return // Ignore if no .gitignore
	}
	content := string(b)
	
	needsUpdate := false
	if !bytes.Contains(b, []byte(".env")) {
		content += "\n.env"
		needsUpdate = true
	}
	if !bytes.Contains(b, []byte(".devdock.local.yml")) {
		content += "\n.devdock.local.yml"
		needsUpdate = true
	}

	if needsUpdate {
		os.WriteFile(gitignorePath, []byte(content), 0644)
	}
}
