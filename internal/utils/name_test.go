package utils_test

import (
	"testing"
	"devdock/internal/utils"
)

func TestNormalizeProjectName(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"My App", "my-app"},
		{"my_app", "my-app"},
		{"some---name", "some-name"},
		{"Trailing- ", "trailing"},
	}
	for _, tt := range tests {
		got := utils.NormalizeProjectName(tt.in)
		if got != tt.want {
			t.Errorf("NormalizeProjectName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestNormalizeDBName(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"my-app", "my_app"},
		{"My App", "my_app"},
		{"some---name", "some_name"},
		{"Trailing_ ", "trailing"},
	}
	for _, tt := range tests {
		got := utils.NormalizeDBName(tt.in)
		if got != tt.want {
			t.Errorf("NormalizeDBName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
