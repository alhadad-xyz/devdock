package process

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

type Runner struct {
	ProjectDir  string
	ProjectName string
	Command     string
}

func NewRunner(projectDir, projectName, command string) *Runner {
	return &Runner{
		ProjectDir:  projectDir,
		ProjectName: projectName,
		Command:     command,
	}
}

func (r *Runner) RunForeground() error {
	cmd := exec.Command("sh", "-c", r.Command)
	cmd.Dir = r.ProjectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nStopping app process...")
		syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
	}()

	err := cmd.Wait()
	fmt.Println("App process stopped.")
	fmt.Println("Docker services are still running. Run 'devdock down' to stop them.")
	return err
}

func (r *Runner) RunDetached() error {
	logsDir := filepath.Join(devdockHome(), "logs")
	os.MkdirAll(logsDir, 0755)
	logPath := filepath.Join(logsDir, r.ProjectName+".app.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("sh", "-c", r.Command)
	cmd.Dir = r.ProjectDir
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return err
	}

	fmt.Printf("Detached app log: %s\n", logPath)
	return WritePID(r.ProjectName, cmd.Process.Pid)
}
