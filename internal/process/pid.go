package process

import (
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

func devdockHome() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".devdock")
}

func getPidPath(projName string) (string, error) {
	pidsDir := filepath.Join(devdockHome(), "pids")
	if err := os.MkdirAll(pidsDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(pidsDir, projName+".pid"), nil
}

func WritePID(projName string, pid int) error {
	path, err := getPidPath(projName)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0644)
}

func ReadPID(projName string) (int, error) {
	path, err := getPidPath(projName)
	if err != nil {
		return 0, err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(b))
}

func ClearPID(projName string) error {
	path, err := getPidPath(projName)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}
