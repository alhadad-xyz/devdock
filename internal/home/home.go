package home

import (
	"fmt"
	"os"
	"path/filepath"
)

// Init initializes the DevDock home directory and its default contents.
func Init() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	devDockDir := filepath.Join(homeDir, ".devdock")
	pidsDir := filepath.Join(devDockDir, "pids")
	logsDir := filepath.Join(devDockDir, "logs")

	dirsToCreate := []string{devDockDir, pidsDir, logsDir}

	for _, dir := range dirsToCreate {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	configFile := filepath.Join(devDockDir, "config.yml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		defaultConfig := `version: "1"
defaults:
  package_manager: pnpm
  editor: code
`
		if err := os.WriteFile(configFile, []byte(defaultConfig), 0644); err != nil {
			return fmt.Errorf("failed to create default config.yml: %w", err)
		}
	}

	return nil
}
