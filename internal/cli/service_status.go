package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"devdock/internal/config"
	"devdock/internal/services"
)

var serviceStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of configured services",
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

		configuredServices := []string{}
		for _, s := range services.All() {
			if isServiceConfigured(cfg, s.Name) {
				configuredServices = append(configuredServices, s.Name)
			}
		}

		if len(configuredServices) == 0 {
			fmt.Println("No services configured.")
			os.Exit(0)
		}

		c := exec.Command("docker", "compose", "ps", "--format", "json")
		c.Dir = projectDir
		out, err := c.Output()
		
		if err != nil || len(out) == 0 {
			fmt.Println("Services are configured but not running. Run `devdock up`.")
			os.Exit(0)
		}

		var dockerServices []map[string]interface{}
		if err := json.Unmarshal(out, &dockerServices); err != nil {
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			for _, line := range lines {
				var s map[string]interface{}
				if json.Unmarshal([]byte(line), &s) == nil {
					dockerServices = append(dockerServices, s)
				}
			}
		}

		if len(dockerServices) == 0 {
			fmt.Println("Services are configured but not running. Run `devdock up`.")
			os.Exit(0)
		}

		dockerMap := make(map[string]map[string]interface{})
		for _, s := range dockerServices {
			name, _ := s["Service"].(string)
			if name == "" {
				name, _ = s["Name"].(string)
			}
			dockerMap[name] = s
		}

		fmt.Printf("%-20s %-15s %-10s %-10s %s\n", "SERVICE", "STATUS", "PORT", "HEALTHY", "WEB URL")
		for _, sName := range configuredServices {
			sData, ok := dockerMap[sName]
			status := "stopped"
			health := "-"
			if ok {
				status, _ = sData["State"].(string)
				h, _ := sData["Health"].(string)
				if h != "" {
					health = h
				}
			}

			def, _ := services.Get(sName)
			primaryPort := getPrimaryPort(cfg, sName, def)
			
			portStr := "-"
			if primaryPort > 0 {
				portStr = fmt.Sprintf("%d", primaryPort)
			}

			webURL := "-"
			if def.WebUI != nil {
				webURL = fmt.Sprintf("http://localhost:%d", getServicePort(cfg, sName, def))
			}

			fmt.Printf("%-20s %-15s %-10s %-10s %s\n", sName, status, portStr, health, webURL)
		}
	},
}

func init() {
	serviceCmd.AddCommand(serviceStatusCmd)
}
