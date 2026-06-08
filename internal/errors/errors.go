package errors

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type ErrorCategory string

const (
	CategoryConfig ErrorCategory = "ConfigError"
	CategoryDocker ErrorCategory = "DockerError"
	CategorySystem ErrorCategory = "SystemError"
	CategoryPort   ErrorCategory = "PortConflict"
)

var exitCodes = map[ErrorCategory]int{
	CategoryConfig: 2,
	CategoryDocker: 3,
	CategoryPort:   4,
	CategorySystem: 1, // Fallback for general errors
}

// AppError represents a centralized application error.
type AppError struct {
	Category ErrorCategory
	What     string
	Why      string
	Fix      string
	Err      error // The underlying cause, not shown to the user unless in debug log
}

// Error implements the error interface. It returns a formatted string for user display.
func (e *AppError) Error() string {
	msg := fmt.Sprintf("✗ %s\n", e.What)
	if e.Why != "" {
		msg += fmt.Sprintf("\n  %s\n", e.Why)
	}
	if e.Fix != "" {
		msg += fmt.Sprintf("\n  Fix: %s\n", e.Fix)
	}
	return msg
}

func (e *AppError) ExitCode() int {
	if code, ok := exitCodes[e.Category]; ok {
		return code
	}
	return 1
}

// Log writes the error details including the stack trace/cause to the error.log.
func (e *AppError) Log() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	logFile := filepath.Join(homeDir, ".devdock", "logs", "error.log")
	
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Format(time.RFC3339)
	logEntry := fmt.Sprintf("[%s] Category: %s | What: %s | Why: %s | Fix: %s", timestamp, e.Category, e.What, e.Why, e.Fix)
	if e.Err != nil {
		logEntry += fmt.Sprintf(" | Cause: %v", e.Err)
	}
	logEntry += "\n"
	
	_, _ = f.WriteString(logEntry)
}

// HandleError prints the AppError to stderr, logs it, and exits with the proper code.
func HandleError(err error) {
	if appErr, ok := err.(*AppError); ok {
		appErr.Log()
		fmt.Fprintln(os.Stderr, appErr.Error())
		os.Exit(appErr.ExitCode())
	} else {
		// Fallback for non-AppError
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
