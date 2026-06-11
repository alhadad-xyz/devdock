package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"devdock/internal/config"
	"devdock/internal/errors"
	"devdock/internal/process"
	"devdock/internal/services"
	"devdock/internal/utils"
)

var openCmd = &cobra.Command{
	Use:   "open [target]",
	Short: "Open the app or a service web UI in the browser",
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

		target := "app"
		if len(args) > 0 {
			target = args[0]
		}

		if target == "app" {
			if cfg.Project.Type == "docker-compose" {
				errors.HandleError(&errors.AppError{
					Category: errors.CategoryConfig,
					What:     "App target is not supported for pure docker-compose projects.",
					Why:      "Docker-compose projects don't have a primary 'app' configured by DevDock.",
					Fix:      "Open specific services instead, e.g., 'devdock open postgres'.",
				})
			}
			
			projName := utils.NormalizeProjectName(cfg.Project.Name)
			pid, err := process.ReadPID(projName)
			isRunning := err == nil && process.IsProcessRunning(pid)
			
			if !isRunning {
				fmt.Println("Warning: App process does not appear to be running. Opening URL anyway.")
			}
			
			port := cfg.App.Port
			if port == 0 {
				errors.HandleError(&errors.AppError{
					Category: errors.CategoryConfig,
					What:     "App port not configured.",
					Why:      "The .devdock.yml file does not specify an app port.",
					Fix:      "Add an 'app.port' value to your .devdock.yml.",
				})
			}
			
			url := fmt.Sprintf("http://localhost:%d", port)
			openBrowser(url)
			return
		}

		// It's a service
		def, ok := services.Get(target)
		if !ok {
				listTargets(cfg)
				errors.HandleError(&errors.AppError{
					Category: errors.CategoryConfig,
					What:     fmt.Sprintf("Unknown target '%s'.", target),
					Why:      "The requested service is not supported by DevDock.",
					Fix:      "Choose from the available targets listed above.",
				})
		}

		if !isServiceConfigured(cfg, target) {
			errors.HandleError(&errors.AppError{
				Category: errors.CategoryConfig,
				What:     fmt.Sprintf("Service '%s' is not enabled.", target),
				Why:      fmt.Sprintf("It exists in DevDock, but hasn't been added to your .devdock.yml config yet."),
				Fix:      fmt.Sprintf("Run `devdock service add %s` to enable it.", target),
			})
		}

		if def.WebUI == nil {
			fmt.Printf("%s does not have a web interface. Connect via your database client at localhost:%d.\n", target, getPrimaryPort(cfg, target, def))
			os.Exit(0)
		}

		port := getServicePort(cfg, target, def)
		url := fmt.Sprintf("http://localhost:%d", port)
		openBrowser(url)
	},
}

func getPrimaryPort(cfg config.Config, target string, def services.ServiceDefinition) int {
	if target == "postgres" && cfg.Services.Postgres != nil && cfg.Services.Postgres.Port != 0 {
		return cfg.Services.Postgres.Port
	}
	if target == "mysql" && cfg.Services.MySQL != nil && cfg.Services.MySQL.Port != 0 {
		return cfg.Services.MySQL.Port
	}
	if target == "redis" && cfg.Services.Redis != nil && cfg.Services.Redis.Port != 0 {
		return cfg.Services.Redis.Port
	}
	// Mailpit primary port is SMTP (1025)
	if target == "mailpit" && cfg.Services.Mailpit != nil && cfg.Services.Mailpit.SMTPPort != 0 {
		return cfg.Services.Mailpit.SMTPPort
	}
	// MinIO primary port is API (9000)
	if target == "minio" && cfg.Services.MinIO != nil && cfg.Services.MinIO.APIPort != 0 {
		return cfg.Services.MinIO.APIPort
	}
	
	for _, p := range def.Ports {
		if p.IsPrimary {
			return p.DefaultPort
		}
	}
	return 0
}

func getServicePort(cfg config.Config, target string, def services.ServiceDefinition) int {
	// Need to resolve the actual configured port for the Web UI.
	if def.WebUI == nil {
		return 0
	}
	
	if target == "mailpit" && cfg.Services.Mailpit != nil && cfg.Services.Mailpit.UIPort != 0 {
		return cfg.Services.Mailpit.UIPort
	}
	if target == "minio" && cfg.Services.MinIO != nil && cfg.Services.MinIO.ConsolePort != 0 {
		return cfg.Services.MinIO.ConsolePort
	}
	
	for _, p := range def.Ports {
		if p.Name == def.WebUI.PortName {
			return p.DefaultPort
		}
	}
	return 0
}

func isServiceConfigured(cfg config.Config, target string) bool {
	if target == "postgres" && cfg.Services.Postgres != nil {
		return cfg.Services.Postgres.Enabled
	}
	if target == "mysql" && cfg.Services.MySQL != nil {
		return cfg.Services.MySQL.Enabled
	}
	if target == "redis" && cfg.Services.Redis != nil {
		return cfg.Services.Redis.Enabled
	}
	if target == "mailpit" && cfg.Services.Mailpit != nil {
		return cfg.Services.Mailpit.Enabled
	}
	if target == "minio" && cfg.Services.MinIO != nil {
		return cfg.Services.MinIO.Enabled
	}
	return false 
}

func listTargets(cfg config.Config) {
	fmt.Println("Available targets to open:")
	if cfg.Project.Type != "docker-compose" {
		fmt.Println("  - app")
	}
	for _, s := range services.All() {
		if isServiceConfigured(cfg, s.Name) {
			if s.WebUI != nil {
				fmt.Printf("  - %s (%s)\n", s.Name, s.WebUI.Label)
			} else {
				fmt.Printf("  - %s (No web interface)\n", s.Name)
			}
		}
	}
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	} else {
		fmt.Printf("Opening %s in your browser...\n", url)
	}
	return err
}

func init() {
	rootCmd.AddCommand(openCmd)
}
