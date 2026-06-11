package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"devdock/internal/config"
	"devdock/internal/errors"
	"devdock/internal/services"
)

var (
	serviceLogsTail  string
	serviceLogsSince string
)

var serviceLogsCmd = &cobra.Command{
	Use:   "logs <service>",
	Short: "View logs for a specific service",
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

		serviceName := args[0]
		
		_, ok := services.Get(serviceName)
		if !ok {
			errors.HandleError(&errors.AppError{
				Category: errors.CategoryConfig,
				What:     fmt.Sprintf("'%s' is not a supported service.", serviceName),
				Why:      "DevDock only supports logs for predefined services.",
				Fix:      "Check the service name and try again.",
			})
		}

		if !isServiceConfigured(cfg, serviceName) {
			errors.HandleError(&errors.AppError{
				Category: errors.CategoryConfig,
				What:     fmt.Sprintf("Service '%s' is not enabled.", serviceName),
				Why:      fmt.Sprintf("It hasn't been added to your .devdock.yml config yet."),
				Fix:      fmt.Sprintf("Run `devdock service add %s` to enable it.", serviceName),
			})
		}

		dockerArgs := []string{"logs", "-f"}
		if serviceLogsTail != "" {
			dockerArgs = append(dockerArgs, "--tail", serviceLogsTail)
		}
		if serviceLogsSince != "" {
			dockerArgs = append(dockerArgs, "--since", serviceLogsSince)
		}
		dockerArgs = append(dockerArgs, serviceName)

		c := exec.Command("docker", append([]string{"compose"}, dockerArgs...)...)
		c.Dir = projectDir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	},
}

func init() {
	serviceLogsCmd.Flags().StringVar(&serviceLogsTail, "tail", "all", "Number of lines to show from the end of the logs")
	serviceLogsCmd.Flags().StringVar(&serviceLogsSince, "since", "", "Show logs since timestamp (e.g. 2013-01-02T13:23:37Z) or relative (e.g. 42m for 42 minutes)")
	serviceCmd.AddCommand(serviceLogsCmd)
}
