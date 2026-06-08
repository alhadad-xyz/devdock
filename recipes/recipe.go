package recipes

import (
	"embed"
	"fmt"

	"gopkg.in/yaml.v3"

	"devdock/internal/config"
)

//go:embed *.yml
var FS embed.FS

func Load(projectType string) (*config.Config, error) {
	b, err := FS.ReadFile(projectType + ".yml")
	if err != nil {
		return nil, fmt.Errorf("no recipe found for project type: %s", projectType)
	}

	var cfg config.Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	// Initialize pointers if nil
	if cfg.Services.Postgres == nil {
		cfg.Services.Postgres = &config.PostgresConfig{}
	}
	if cfg.Services.MySQL == nil {
		cfg.Services.MySQL = &config.MySQLConfig{}
	}
	if cfg.Services.Redis == nil {
		cfg.Services.Redis = &config.RedisConfig{}
	}

	return &cfg, nil
}
