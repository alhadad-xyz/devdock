package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"devdock/internal/config"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system requirements and project configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running DevDock Doctor...")
		fmt.Println()

		// DD-057 Check Docker installed
		err := checkCommand("docker", "--version")
		printCheck("Docker installed", err)

		// DD-058 Check Docker daemon running
		err = checkCommand("docker", "info")
		printCheck("Docker daemon running", err)

		// DD-059 Check Docker Compose v2
		err = checkCommand("docker", "compose", "version")
		printCheck("Docker Compose v2", err)

		// DD-060, DD-061 .devdock.yml handling
		projectDir, _ := ResolveProjectRoot()
		configPath := filepath.Join(projectDir, ".devdock.yml")

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// DD-061 show helpful message but do not crash
			fmt.Println("\nℹ No .devdock.yml found in current directory. Run `devdock init` to create one.")
		} else {
			cfg, err := config.Load(projectDir)
			if err != nil {
				printCheck("Project configuration", fmt.Errorf("invalid .devdock.yml: %v", err))
			} else {
				printCheck("Project configuration (.devdock.yml is valid)", nil)
				
				if cfg.Project.Type == "laravel" {
					err = checkCommand("php", "-v")
					printCheck("PHP runtime", err)
				} else if cfg.Project.Type == "nextjs" || cfg.Project.Type == "express" {
					err = checkCommand("node", "-v")
					printCheck("Node.js runtime", err)
				} else if cfg.Project.Type == "fiber" {
					err = checkCommand("go", "version")
					printCheck("Go runtime", err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func checkCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func printCheck(name string, err error) {
	if err == nil {
		fmt.Printf("✔ %s\n", name)
	} else {
		fmt.Printf("✗ %s (%v)\n", name, err)
	}
}
