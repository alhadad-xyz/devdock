package env

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func MergeSafe(projectDir string, vars map[string]string) error {
	if len(vars) == 0 {
		return nil
	}

	envPath := filepath.Join(projectDir, ".env")
	
	existingKeys := make(map[string]bool)
	b, err := os.ReadFile(envPath)
	if err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(b))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) > 0 {
				existingKeys[strings.TrimSpace(parts[0])] = true
			}
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	var toAppend []string
	for k, v := range vars {
		if !existingKeys[k] {
			toAppend = append(toAppend, fmt.Sprintf("%s=%s", k, v))
		}
	}

	if len(toAppend) == 0 {
		return nil
	}

	f, err := os.OpenFile(envPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if len(b) > 0 && !bytes.HasSuffix(b, []byte("\n")) {
		f.WriteString("\n")
	}

	for _, line := range toAppend {
		f.WriteString(line + "\n")
	}

	return nil
}
