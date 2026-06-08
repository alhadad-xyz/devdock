package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"devdock/internal/config"
	"devdock/internal/process"
	"devdock/internal/utils"
)

var (
	logsTail  string
	logsSince string
)

var logsCmd = &cobra.Command{
	Use:   "logs [service]",
	Short: "View logs from services",
	Run: func(cmd *cobra.Command, args []string) {
		projectDir, _ := ResolveProjectRoot()
		configPath := filepath.Join(projectDir, ".devdock.yml")

		b, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Println("DevDock is not initialized.")
			os.Exit(1)
		}

		var cfg config.Config
		yaml.Unmarshal(b, &cfg)

		service := ""
		if len(args) > 0 {
			service = args[0]
		}

		if service == "app" {
			if cfg.Project.Type == "docker-compose" {
				fmt.Println("This is a docker-compose project. The host app process is not managed by DevDock.")
				os.Exit(1)
			}
			
			projName := utils.NormalizeProjectName(cfg.Project.Name)
			pid, err := process.ReadPID(projName)
			if err != nil || !process.IsProcessRunning(pid) {
				fmt.Println("App process is not running in detached mode.")
				os.Exit(1)
			}

			home, _ := os.UserHomeDir()
			logPath := filepath.Join(home, ".devdock", "logs", projName+".app.log")
			
			if _, err := os.Stat(logPath); os.IsNotExist(err) {
				fmt.Println("App logs are unavailable.")
				os.Exit(1)
			}

			tailArgs := []string{"-f"}
			if logsTail != "" {
				tailArgs = append(tailArgs, "-n", logsTail)
			}
			tailArgs = append(tailArgs, logPath)
			c := exec.Command("tail", tailArgs...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Run()
			return
		}

		dockerArgs := []string{"logs", "-f"}
		if logsTail != "" {
			dockerArgs = append(dockerArgs, "--tail", logsTail)
		}
		if logsSince != "" {
			dockerArgs = append(dockerArgs, "--since", logsSince)
		}
		if service != "" {
			dockerArgs = append(dockerArgs, service)
		}

		c := exec.Command("docker", append([]string{"compose"}, dockerArgs...)...)
		c.Dir = projectDir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	},
}

func init() {
	logsCmd.Flags().StringVarP(&logsTail, "tail", "t", "", "Number of lines to show from the end of the logs")
	logsCmd.Flags().StringVar(&logsSince, "since", "", "Show logs since timestamp (e.g. 2013-01-02T13:23:37Z) or relative (e.g. 42m for 42 minutes)")
	rootCmd.AddCommand(logsCmd)
}
