package config

import (
	"testing"
)

func TestNonExistentConfig(t *testing.T) {
	// Non-existent config file should not cause a failure as its not needed.
	c, err := New("non-existent-file.json")
	if err != nil {
		t.Fatalf("expected test to succeed but it failed")
		return
	}

	expected := Config{}
	if c != expected {
		t.Fatalf("expected to receive empty config but got: %v", c)
	}
}
