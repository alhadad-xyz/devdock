package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"devdock/internal/config"
	"devdock/internal/services"
)

var restartCmd = &cobra.Command{
	Use:   "restart [service]",
	Short: "Restart a specific service container",
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

		if len(args) == 0 {
			fmt.Println("Specify a service to restart. Example: `devdock restart postgres`")
			os.Exit(1)
		}

		serviceName := args[0]

		if serviceName == "app" {
			fmt.Println("Restarting the app is not supported. Use Ctrl+C then `devdock up`.")
			os.Exit(0)
		}

		_, ok := services.Get(serviceName)
		if !ok || !isServiceConfigured(cfg, serviceName) {
			fmt.Printf("Unknown or unconfigured service: '%s'\n", serviceName)
			fmt.Println("Configured services:")
			for _, s := range services.All() {
				if isServiceConfigured(cfg, s.Name) {
					fmt.Printf("  - %s\n", s.Name)
				}
			}
			os.Exit(1)
		}

		fmt.Printf("Restarting %s...\n", serviceName)
		c := exec.Command("docker", "compose", "restart", serviceName)
		c.Dir = projectDir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			fmt.Printf("Failed to restart %s: %v\n", serviceName, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
