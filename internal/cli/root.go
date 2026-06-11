package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Version metadata
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

var (
	quietFlag   bool
	jsonFlag    bool
	projectFlag string
)

var rootCmd = &cobra.Command{
	Use:           "devdock",
	Short:         "DevDock is a local development environment manager",
	Long:          `DevDock simplifies local development for Laravel and Next.js applications using Docker Compose.`,
	SilenceErrors: true, // Handled manually in Execute
	SilenceUsage:  true, // Prevent showing usage on every error
	Version:       Version,
	Run: func(cmd *cobra.Command, args []string) {
		if Version == "dev" {
			fmt.Printf("Welcome to DevDock v%s\n\n", Version)
		} else {
			fmt.Printf("Welcome to DevDock v%s (commit: %s, built: %s)\n\n", Version, Commit, BuildDate)
		}
		fmt.Println("Run `devdock --help` to see available commands.")
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "Suppress output")
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().StringVar(&projectFlag, "project", "", "Project root path (default: current directory)")
}

// ResolveProjectRoot resolves the project path from --project or current directory.
func ResolveProjectRoot() (string, error) {
	if projectFlag != "" {
		return filepath.Abs(projectFlag)
	}
	return os.Getwd()
}

// Execute executes the root command.
func Execute() error {
	if Version == "dev" {
		rootCmd.SetVersionTemplate("devdock version dev\n")
	} else {
		rootCmd.SetVersionTemplate(fmt.Sprintf("devdock version %s (commit: %s, built: %s)\n", Version, Commit, BuildDate))
	}

	err := rootCmd.Execute()
	if err != nil {
		if strings.HasPrefix(err.Error(), "unknown command") {
			parts := strings.Split(err.Error(), "\"")
			cmdName := "unknown"
			if len(parts) >= 3 {
				cmdName = parts[1]
			}
			fmt.Fprintf(os.Stderr, "'devdock %s' is not available in this version.\n\n", cmdName)
			fmt.Fprintln(os.Stderr, "Run `devdock --help` to see available commands.")
			os.Exit(1)
		}
		// For other errors
		return err
	}
	return nil
}
