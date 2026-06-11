package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"devdock/internal/config"
	"devdock/internal/process"
	"devdock/internal/utils"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of the DevDock environment",
	Run: func(cmd *cobra.Command, args []string) {
		projectDir, _ := ResolveProjectRoot()
		configPath := filepath.Join(projectDir, ".devdock.yml")

		b, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Println("DevDock is not initialized. Run 'devdock init' first.")
			os.Exit(1)
		}

		var cfg config.Config
		yaml.Unmarshal(b, &cfg)

		fmt.Printf("Project: %s\n", cfg.Project.Name)
		fmt.Printf("Path:    %s\n", projectDir)
		fmt.Printf("Type:    %s\n", cfg.Project.Type)

		fmt.Println("\nDocker Services:")
		
		c := exec.Command("docker", "compose", "ps", "--format", "json")
		c.Dir = projectDir
		out, err := c.Output()
		
		if err == nil && len(out) > 0 {
			var services []map[string]interface{}
			if err := json.Unmarshal(out, &services); err != nil {
				lines := strings.Split(strings.TrimSpace(string(out)), "\n")
				for _, line := range lines {
					var s map[string]interface{}
					if json.Unmarshal([]byte(line), &s) == nil {
						services = append(services, s)
					}
				}
			}

			if len(services) > 0 {
				fmt.Printf("  %-20s %-15s %s\n", "NAME", "STATE", "HEALTH")
				for _, s := range services {
					name, _ := s["Service"].(string)
					if name == "" {
						name, _ = s["Name"].(string)
					}
					state, _ := s["State"].(string)
					health, _ := s["Health"].(string)
					if health == "" {
						health = "-"
					}
					fmt.Printf("  %-20s %-15s %s\n", name, state, health)
				}
			} else {
				fmt.Println("  No running services.")
			}
		} else {
			fmt.Println("  No running services.")
		}

		if cfg.Project.Type != "docker-compose" {
			fmt.Println("\nApp Process:")
			projName := utils.NormalizeProjectName(cfg.Project.Name)
			pid, err := process.ReadPID(projName)
			if err == nil && process.IsProcessRunning(pid) {
				fmt.Printf("  Running (Detached PID: %d)\n", pid)
			} else {
				fmt.Println("  Stopped (or running in foreground)")
			}

			if cfg.App.Port != 0 {
				fmt.Printf("\n  URL: http://localhost:%d\n", cfg.App.Port)
			}
		}

		hasDB := false
		if cfg.Services.Postgres != nil && cfg.Services.Postgres.Enabled {
			hasDB = true
			fmt.Printf("\nPostgres: postgresql://postgres:postgres@127.0.0.1:%d/postgres\n", cfg.Services.Postgres.Port)
		}
		if cfg.Services.MySQL != nil && cfg.Services.MySQL.Enabled {
			if !hasDB {
				fmt.Println()
				hasDB = true
			}
			dbName := utils.NormalizeDBName(cfg.Project.Name)
			fmt.Printf("MySQL: mysql://root:root@127.0.0.1:%d/%s\n", cfg.Services.MySQL.Port, dbName)
		}
		if cfg.Services.Redis != nil && cfg.Services.Redis.Enabled {
			if !hasDB {
				fmt.Println()
			}
			fmt.Printf("Redis: redis://127.0.0.1:%d\n", cfg.Services.Redis.Port)
		}
		if cfg.Services.Mailpit != nil && cfg.Services.Mailpit.Enabled {
			fmt.Printf("Mailpit Web UI: http://127.0.0.1:%d\n", cfg.Services.Mailpit.UIPort)
		}
		if cfg.Services.MinIO != nil && cfg.Services.MinIO.Enabled {
			fmt.Printf("MinIO Console: http://127.0.0.1:%d\n", cfg.Services.MinIO.ConsolePort)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
