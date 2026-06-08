package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"

	"devdock/internal/errors"
)

var projectNameRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

// Load reads and validates .devdock.yml from the project directory.
func Load(projectDir string) (*Config, error) {
	configPath := filepath.Join(projectDir, ".devdock.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &errors.AppError{
				Category: errors.CategoryConfig,
				What:     "Config file missing",
				Why:      fmt.Sprintf("Could not find .devdock.yml in %s", projectDir),
				Fix:      "Run `devdock init` to create a new project configuration.",
				Err:      err,
			}
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, &errors.AppError{
			Category: errors.CategoryConfig,
			What:     "Invalid config file",
			Why:      fmt.Sprintf(".devdock.yml contains invalid YAML: %v", err),
			Fix:      "Fix the YAML syntax errors in .devdock.yml.",
			Err:      err,
		}
	}

	// Validation
	if cfg.Version != "1" {
		return nil, &errors.AppError{
			Category: errors.CategoryConfig,
			What:     "Unsupported config version",
			Why:      fmt.Sprintf("Expected version '1', got '%s'", cfg.Version),
			Fix:      "Change 'version: \"1\"' in .devdock.yml.",
		}
	}

	if cfg.Project.Name == "" {
		return nil, &errors.AppError{
			Category: errors.CategoryConfig,
			What:     "Missing project name",
			Why:      "project.name is required",
			Fix:      "Add 'name' under 'project' in .devdock.yml.",
		}
	}

	if !projectNameRegex.MatchString(cfg.Project.Name) {
		return nil, &errors.AppError{
			Category: errors.CategoryConfig,
			What:     "Invalid project name",
			Why:      "project.name must match ^[a-z0-9-]+$",
			Fix:      "Change project.name to use only lowercase letters, numbers, and hyphens.",
		}
	}

	switch cfg.Project.Type {
	case "laravel", "nextjs":
		if cfg.App.Command == "" {
			return nil, &errors.AppError{
				Category: errors.CategoryConfig,
				What:     "Missing app command",
				Why:      fmt.Sprintf("app.command is required for %s projects", cfg.Project.Type),
				Fix:      "Add 'command' under 'app' in .devdock.yml.",
			}
		}
		if cfg.App.Port == 0 {
			return nil, &errors.AppError{
				Category: errors.CategoryConfig,
				What:     "Missing app port",
				Why:      fmt.Sprintf("app.port is required for %s projects", cfg.Project.Type),
				Fix:      "Add 'port' under 'app' in .devdock.yml.",
			}
		}
	case "docker-compose":
		// Optional and ignored.
	default:
		return nil, &errors.AppError{
			Category: errors.CategoryConfig,
			What:     "Invalid project type",
			Why:      fmt.Sprintf("'%s' is not a supported project type", cfg.Project.Type),
			Fix:      "Change project.type to one of: laravel, nextjs, docker-compose.",
		}
	}

	if cfg.App.RunMode == "container" {
		return nil, &errors.AppError{
			Category: errors.CategoryConfig,
			What:     "Unsupported run mode",
			Why:      "app.run_mode: container is not supported yet in v0.1.",
			Fix:      "Remove run_mode or change it to 'host' in .devdock.yml.",
		}
	} else if cfg.App.RunMode != "" && cfg.App.RunMode != "host" {
		return nil, &errors.AppError{
			Category: errors.CategoryConfig,
			What:     "Invalid run mode",
			Why:      fmt.Sprintf("run_mode '%s' is not supported", cfg.App.RunMode),
			Fix:      "Change run_mode to 'host' or remove it.",
		}
	}

	// Validate ports
	usedPorts := make(map[int]string)

	if cfg.Services.Postgres != nil && cfg.Services.Postgres.Enabled && cfg.Services.Postgres.Port != 0 {
		if err := validatePort(cfg.Services.Postgres.Port, "services.postgres", usedPorts); err != nil {
			return nil, err
		}
	}
	if cfg.Services.MySQL != nil && cfg.Services.MySQL.Enabled && cfg.Services.MySQL.Port != 0 {
		if err := validatePort(cfg.Services.MySQL.Port, "services.mysql", usedPorts); err != nil {
			return nil, err
		}
	}
	if cfg.Services.Redis != nil && cfg.Services.Redis.Enabled && cfg.Services.Redis.Port != 0 {
		if err := validatePort(cfg.Services.Redis.Port, "services.redis", usedPorts); err != nil {
			return nil, err
		}
	}

	return &cfg, nil
}

func validatePort(port int, field string, usedPorts map[int]string) error {
	if port < 1 || port > 65535 {
		return &errors.AppError{
			Category: errors.CategoryConfig,
			What:     "Invalid port",
			Why:      fmt.Sprintf("Port %d in %s is out of valid range (1-65535)", port, field),
			Fix:      fmt.Sprintf("Change %s.port to a valid TCP port.", field),
		}
	}
	if prevField, ok := usedPorts[port]; ok {
		return &errors.AppError{
			Category: errors.CategoryConfig,
			What:     "Duplicate port",
			Why:      fmt.Sprintf("Port %d is configured for both %s and %s", port, prevField, field),
			Fix:      "Assign a unique port for each service in .devdock.yml.",
		}
	}
	usedPorts[port] = field
	return nil
}
