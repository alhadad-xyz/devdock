package ports

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

type ConflictError struct {
	Port    int
	Service string
	PID     string
	Process string
	Fix     string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("✗ Port %d is already in use\n\n  PID %s (%s) is using this port.\n\n  Fix: %s", e.Port, e.PID, e.Process, e.Fix)
}

func CheckPort(port int, serviceName string, fixMsg string) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		// Port is in use
		pid, proc := getProcessUsingPort(port)
		return &ConflictError{
			Port:    port,
			Service: serviceName,
			PID:     pid,
			Process: proc,
			Fix:     fixMsg,
		}
	}
	l.Close()
	return nil
}

func getProcessUsingPort(port int) (string, string) {
	cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port), "-t")
	out, err := cmd.Output()
	if err != nil || len(out) == 0 {
		return "unknown", "unknown"
	}
	pids := strings.Split(strings.TrimSpace(string(out)), "\n")
	pid := pids[0]

	cmd2 := exec.Command("ps", "-p", pid, "-o", "comm=")
	out2, err2 := cmd2.Output()
	proc := "unknown"
	if err2 == nil && len(out2) > 0 {
		proc = strings.TrimSpace(string(out2))
		parts := strings.Split(proc, "/")
		proc = parts[len(parts)-1]
	}

	return pid, proc
}
