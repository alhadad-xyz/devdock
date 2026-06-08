package docker

import (
	"os"
	"os/exec"
)

func runCompose(projectDir string, args ...string) error {
	cmd := exec.Command("docker", append([]string{"compose"}, args...)...)
	cmd.Dir = projectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Up(projectDir string, build bool) error {
	args := []string{"up", "-d"}
	if build {
		args = append(args, "--build")
	}
	return runCompose(projectDir, args...)
}

func Down(projectDir string, volumes bool) error {
	args := []string{"down"}
	if volumes {
		args = append(args, "--volumes")
	}
	return runCompose(projectDir, args...)
}

func Ps(projectDir string) error {
	return runCompose(projectDir, "ps")
}

func Logs(projectDir string, service string) error {
	args := []string{"logs", "-f"}
	if service != "" {
		args = append(args, service)
	}
	return runCompose(projectDir, args...)
}
