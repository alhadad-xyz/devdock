package services

import (
	"testing"
)

func TestRegistryGet(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"postgres", true},
		{"mysql", true},
		{"redis", true},
		{"mailpit", true},
		{"minio", true},
		{"mongodb", false},
	}

	for _, tt := range tests {
		_, ok := Get(tt.name)
		if ok != tt.expected {
			t.Errorf("Get(%q) expected %v, got %v", tt.name, tt.expected, ok)
		}
	}
}

func TestRegistryAll(t *testing.T) {
	all := All()
	if len(all) != 5 {
		t.Errorf("Expected 5 registered services, got %d", len(all))
	}
}
