package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"devdock/internal/compose"
	"devdock/internal/config"
	"devdock/internal/errors"
	"devdock/internal/services"
)

var serviceRemoveCmd = &cobra.Command{
	Use:   "remove <service>",
	Short: "Remove a service from your project",
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
				Fix:      "Please edit your compose.yml file directly.",
			})
		}

		serviceName := args[0]
		if !isServiceConfigured(cfg, serviceName) {
			fmt.Printf("ℹ %s is not configured in this project.\n", serviceName)
			os.Exit(0)
		}

		def, ok := services.Get(serviceName)
		if ok && def.HasVolume {
			fmt.Printf("Warning: %s uses a persistent volume. Removing it will not delete the data automatically.\n", serviceName)
			fmt.Println("To delete data, you must run `devdock down --volumes` before or after removal.")
		}

		confirm := false
		survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("Remove %s from this project?", serviceName),
			Default: false,
		}, &confirm)

		if !confirm {
			fmt.Println("No changes made.")
			os.Exit(0)
		}

		// Remove from config
		if serviceName == "postgres" && cfg.Services.Postgres != nil {
			cfg.Services.Postgres.Enabled = false
		} else if serviceName == "mysql" && cfg.Services.MySQL != nil {
			cfg.Services.MySQL.Enabled = false
		} else if serviceName == "redis" && cfg.Services.Redis != nil {
			cfg.Services.Redis.Enabled = false
		} else if serviceName == "mailpit" && cfg.Services.Mailpit != nil {
			cfg.Services.Mailpit.Enabled = false
		} else if serviceName == "minio" && cfg.Services.MinIO != nil {
			cfg.Services.MinIO.Enabled = false
		}

		tmpPath := configPath + ".tmp"
		newB, _ := yaml.Marshal(&cfg)
		os.WriteFile(tmpPath, newB, 0644)
		os.Rename(tmpPath, configPath)

		compose.Generate(projectDir, &cfg)

		fmt.Println("Run `devdock down && devdock up` to apply.")
	},
}

func init() {
	serviceCmd.AddCommand(serviceRemoveCmd)
}
