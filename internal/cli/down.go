package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"devdock/internal/config"
	"devdock/internal/docker"
	"devdock/internal/process"
	"devdock/internal/utils"
)

var (
	downVolumes bool
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop the DevDock environment",
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

		projName := utils.NormalizeProjectName(cfg.Project.Name)

		if pid, err := process.ReadPID(projName); err == nil {
			if process.IsProcessRunning(pid) {
				fmt.Printf("Stopping app process (PID %d)...\n", pid)
				p, _ := os.FindProcess(pid)
				// Send SIGINT or SIGTERM. We used SIGINT in foreground so let's use SIGTERM here.
				p.Signal(syscall.SIGTERM)
				
				for i := 0; i < 10; i++ {
					if !process.IsProcessRunning(pid) {
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
			}
			process.ClearPID(projName)
		}

		if downVolumes {
			confirm := false
			err := survey.AskOne(&survey.Confirm{
				Message: "This will destroy all database volumes and data. Are you sure?",
				Default: false,
			}, &confirm)
			if err != nil || !confirm {
				fmt.Println("Aborted volume destruction. Services will just be stopped.")
				downVolumes = false
			}
		}

		fmt.Println("Stopping Docker services...")
		docker.Down(projectDir, downVolumes)
		fmt.Println("DevDock environment stopped.")
	},
}

func init() {
	downCmd.Flags().BoolVarP(&downVolumes, "volumes", "v", false, "Remove named volumes (destroys data)")
	rootCmd.AddCommand(downCmd)
}
