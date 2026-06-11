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
)

var runCmd = &cobra.Command{
	Use:   "run [command]",
	Short: "Run a command defined in your .devdock.yml",
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
			listCommands(cfg)
			os.Exit(0)
		}

		commandName := args[0]
		
		if cfg.Commands == nil || len(cfg.Commands) == 0 {
			errors.HandleError(&errors.AppError{
				Category: errors.CategoryConfig,
				What:     "No commands defined.",
				Why:      "Your .devdock.yml file does not contain a 'commands' block.",
				Fix:      "Add commands to .devdock.yml (e.g. 'dev: npm run dev').",
			})
		}
		
		commandString, ok := cfg.Commands[commandName]
		if !ok {
			listCommands(cfg)
			errors.HandleError(&errors.AppError{
				Category: errors.CategoryConfig,
				What:     fmt.Sprintf("Command '%s' is not defined.", commandName),
				Why:      "The command you specified does not exist in your .devdock.yml file.",
				Fix:      "Choose one of the available commands listed above, or add it to .devdock.yml.",
			})
		}

		c := exec.Command("/bin/sh", "-c", commandString)
		c.Dir = projectDir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				os.Exit(exitError.ExitCode())
			}
			fmt.Printf("Failed to run command: %v\n", err)
			os.Exit(1)
		}
	},
}

func listCommands(cfg config.Config) {
	fmt.Println("Available commands:")
	if len(cfg.Commands) == 0 {
		fmt.Println("  (None defined)")
		return
	}
	
	// Print in a nice table
	for name, cmdStr := range cfg.Commands {
		fmt.Printf("  %-15s %s\n", name, cmdStr)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
}
