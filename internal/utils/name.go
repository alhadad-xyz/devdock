package utils

import (
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// NormalizeProjectName converts names to lowercase kebab-case (e.g. "My App" -> "my-app", "my_app" -> "my-app")
func NormalizeProjectName(name string) string {
	s := strings.ToLower(name)
	s = nonAlphanumericRegex.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// NormalizeDBName converts names to lowercase snake_case (e.g. "my-app" -> "my_app")
func NormalizeDBName(name string) string {
	s := strings.ToLower(name)
	s = nonAlphanumericRegex.ReplaceAllString(s, "_")
	return strings.Trim(s, "_")
}
