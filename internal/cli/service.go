package cli

import (
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage services in your DevDock project",
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
