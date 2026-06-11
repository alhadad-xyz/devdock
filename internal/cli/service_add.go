package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"devdock/internal/compose"
	"devdock/internal/config"
	"devdock/internal/env"
	"devdock/internal/errors"
	"devdock/internal/services"
)

var serviceAddCmd = &cobra.Command{
	Use:   "add <service>",
	Short: "Add a service to your project",
	Args:  cobra.ExactArgs(1),
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

		if cfg.Project.Type == "docker-compose" {
			errors.HandleError(&errors.AppError{
				Category: errors.CategoryConfig,
				What:     "Service management is disabled for pure docker-compose projects.",
				Why:      "DevDock does not manage the compose.yml directly for this project type.",
				Fix:      "Please edit your compose.yml file directly to add services.",
			})
		}

		serviceName := args[0]
		_, ok := services.Get(serviceName)
		if !ok {
			fmt.Println("Supported services:")
			for _, s := range services.All() {
				fmt.Printf("  - %s\n", s.Name)
			}
			errors.HandleError(&errors.AppError{
				Category: errors.CategoryConfig,
				What:     fmt.Sprintf("'%s' is not a supported service.", serviceName),
				Why:      "DevDock only supports adding predefined services via this CLI.",
				Fix:      "Choose from the supported services listed above.",
			})
		}

		if isServiceConfigured(cfg, serviceName) {
			fmt.Printf("%s is already configured in this project.\n", serviceName)
			os.Exit(0)
		}

		// Add service to config
		if serviceName == "postgres" {
			if cfg.Services.Postgres == nil {
				cfg.Services.Postgres = &config.PostgresConfig{}
			}
			cfg.Services.Postgres.Enabled = true
			cfg.Services.Postgres.Port = 5432
			cfg.Services.Postgres.Version = "15"
		} else if serviceName == "mysql" {
			if cfg.Services.MySQL == nil {
				cfg.Services.MySQL = &config.MySQLConfig{}
			}
			cfg.Services.MySQL.Enabled = true
			cfg.Services.MySQL.Port = 3306
			cfg.Services.MySQL.Version = "8.0"
		} else if serviceName == "redis" {
			if cfg.Services.Redis == nil {
				cfg.Services.Redis = &config.RedisConfig{}
			}
			cfg.Services.Redis.Enabled = true
			cfg.Services.Redis.Port = 6379
			cfg.Services.Redis.Version = "7"
		} else if serviceName == "mailpit" {
			if cfg.Services.Mailpit == nil {
				cfg.Services.Mailpit = &config.MailpitConfig{}
			}
			cfg.Services.Mailpit.Enabled = true
			cfg.Services.Mailpit.Version = "v1.21"
			cfg.Services.Mailpit.SMTPPort = 1025
			cfg.Services.Mailpit.UIPort = 8025
		} else if serviceName == "minio" {
			if cfg.Services.MinIO == nil {
				cfg.Services.MinIO = &config.MinIOConfig{}
			}
			cfg.Services.MinIO.Enabled = true
			cfg.Services.MinIO.Version = "RELEASE.2024-11-07T00-52-20Z"
			cfg.Services.MinIO.APIPort = 9000
			cfg.Services.MinIO.ConsolePort = 9001
		}

		// Write config atomically
		tmpPath := configPath + ".tmp"
		newB, _ := yaml.Marshal(&cfg)
		os.WriteFile(tmpPath, newB, 0644)
		os.Rename(tmpPath, configPath)

		// Regenerate compose
		compose.Generate(projectDir, &cfg)

		vars := env.GetServiceVars(cfg.Project.Name, cfg.Project.Type, serviceName, &cfg)
		if len(vars) > 0 {
			fmt.Printf("\nEnvironment variables for %s:\n", cfg.Project.Type)
			for k, v := range vars {
				fmt.Printf("%s=%s\n", k, v)
			}

			confirm := true
			err := survey.AskOne(&survey.Confirm{
				Message: fmt.Sprintf("Append these to .env?"),
				Default: true,
			}, &confirm)

			if err == nil && confirm {
				env.MergeSafe(projectDir, vars)
			}
		}

		// Check if running
		c := exec.Command("docker", "compose", "ps", "-q")
		c.Dir = projectDir
		out, _ := c.Output()
		if len(out) > 0 {
			fmt.Println("\nℹ Notice: Docker services are currently running.")
			fmt.Println("  Run `devdock down && devdock up` to start your new service.")
		} else {
			fmt.Println("\nRun `devdock up` to start your environment.")
		}
	},
}

func init() {
	serviceCmd.AddCommand(serviceAddCmd)
}
