package docker

import (
	"fmt"
	"os/exec"
	"strings"
)

func CheckPrerequisites() error {
	// Check docker installed
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker is not installed. Please install Docker Desktop or Docker Engine")
	}

	// Check daemon running
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker daemon is not running. Please start Docker")
	}

	// Check compose v2
	cmd = exec.Command("docker", "compose", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker compose is not available")
	}
	
	if !strings.Contains(string(out), "version v2") {
		// Just a warning or return error?
		// We can return nil but ideally we need v2
	}

	return nil
}

func CheckNode() error {
	if _, err := exec.LookPath("node"); err != nil {
		return fmt.Errorf("node is not installed")
	}
	if _, err := exec.LookPath("npm"); err != nil {
		return fmt.Errorf("npm is not installed")
	}
	return nil
}

func CheckPHP() error {
	if _, err := exec.LookPath("php"); err != nil {
		return fmt.Errorf("php is not installed")
	}
	return nil
}

func CheckGo() error {
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("go is not installed")
	}
	return nil
}
