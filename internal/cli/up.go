package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"devdock/internal/compose"
	"devdock/internal/config"
	"devdock/internal/docker"
	"devdock/internal/ports"
	"devdock/internal/process"
	"devdock/internal/utils"
)

var (
	upDetach     bool
	upBuild      bool
	upSkipChecks bool
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the DevDock environment",
	Run: func(cmd *cobra.Command, args []string) {
		projectDir, _ := ResolveProjectRoot()
		configPath := filepath.Join(projectDir, ".devdock.yml")

		b, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Println("DevDock is not initialized. Run 'devdock init' first.")
			os.Exit(1)
		}

		var cfg config.Config
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			fmt.Printf("Error reading config: %v\n", err)
			os.Exit(1)
		}

		if !upSkipChecks {
			if err := docker.CheckPrerequisites(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if cfg.Project.Type == "laravel" {
				if err := docker.CheckPHP(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else if cfg.Project.Type == "nextjs" {
				if err := docker.CheckNode(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			if cfg.App.Port != 0 {
				if err := ports.CheckPort(cfg.App.Port, "app", fmt.Sprintf("Change app.port in .devdock.yml to a different port, then run 'devdock up' again.")); err != nil {
					fmt.Println(err)
					os.Exit(4)
				}
			}

			if cfg.Services.Postgres != nil && cfg.Services.Postgres.Enabled {
				if err := ports.CheckPort(cfg.Services.Postgres.Port, "postgres", fmt.Sprintf("Change services.postgres.port in .devdock.yml to a different port, then run 'devdock up' again.")); err != nil {
					fmt.Println(err)
					os.Exit(4)
				}
			}

			if cfg.Services.MySQL != nil && cfg.Services.MySQL.Enabled {
				if err := ports.CheckPort(cfg.Services.MySQL.Port, "mysql", fmt.Sprintf("Change services.mysql.port in .devdock.yml to a different port, then run 'devdock up' again.")); err != nil {
					fmt.Println(err)
					os.Exit(4)
				}
			}

			if cfg.Services.Redis != nil && cfg.Services.Redis.Enabled {
				if err := ports.CheckPort(cfg.Services.Redis.Port, "redis", fmt.Sprintf("Change services.redis.port in .devdock.yml to a different port, then run 'devdock up' again.")); err != nil {
					fmt.Println(err)
					os.Exit(4)
				}
			}
		}

		projName := utils.NormalizeProjectName(cfg.Project.Name)
		if pid, err := process.ReadPID(projName); err == nil {
			if process.IsProcessRunning(pid) {
				fmt.Printf("App process is already running in detached mode (PID %d).\n", pid)
				fmt.Println("Run 'devdock down' first.")
				os.Exit(1)
			} else {
				process.ClearPID(projName)
			}
		}

		if cfg.Project.Type != "docker-compose" {
			compose.Generate(projectDir, &cfg)
		}

		fmt.Println("Starting Docker services...")
		if err := docker.Up(projectDir, upBuild); err != nil {
			fmt.Printf("Error starting services. A service may be unhealthy or failed to start.\n")
			fmt.Printf("Run 'devdock logs <service>' or 'docker compose logs' to investigate.\n")
			os.Exit(5)
		}

		if cfg.Project.Type == "docker-compose" {
			fmt.Println("\nDocker Compose services started.")
			return
		}

		if cfg.App.Command != "" {
			runner := process.NewRunner(projectDir, projName, cfg.App.Command)
			
			if upDetach {
				fmt.Println("Starting app process in background...")
				if err := runner.RunDetached(); err != nil {
					fmt.Printf("Error starting app: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("App running at http://localhost:%d\n", cfg.App.Port)
			} else {
				fmt.Printf("App running at http://localhost:%d\n", cfg.App.Port)
				fmt.Println("Starting app process...")
				runner.RunForeground()
			}
		}
	},
}

func init() {
	upCmd.Flags().BoolVarP(&upDetach, "detach", "d", false, "Run in background")
	upCmd.Flags().BoolVar(&upBuild, "build", false, "Build images before starting containers")
	upCmd.Flags().BoolVar(&upSkipChecks, "skip-checks", false, "Skip prerequisite and port checks")
	rootCmd.AddCommand(upCmd)
}
