package parser

import "testing"

func TestValidateEndpoint(t *testing.T) {
	// Empty endpoint should be set to default value
	if endpoint, err := validateEndpoint(""); err != nil {
		t.Errorf("Received unexpected error when validating empty endpoint")
	} else if endpoint != DefaultEndpoint {
		t.Errorf("Empty endpoint should produce default endpoint")
	}

	// Endpoint without leading slash should fail
	if _, err := validateEndpoint("testing"); err == nil {
		t.Errorf("Expected to receive error when validating endpoint without leading slash")
	}

	// Valid endpoint should produce itself
	if endpoint, err := validateEndpoint("/users/apicheck"); err != nil {
		t.Errorf("Received unexpected error when validating endpoint")
	} else if endpoint != "/users/apicheck" {
		t.Errorf("Did not receive expected endpoint")
	}
}

func TestValidateHostname(t *testing.T) {
	// Empty hostname should be an error
	if _, err := validateHostname(""); err == nil {
		t.Errorf("Expected to receive error when validating empty hostname")
	}

	// Valid hostname should succeed
	if endpoint, err := validateHostname("http://localhost:3000"); err != nil {
		t.Errorf("Received unexpected error when validating hostname")
	} else if endpoint != "http://localhost:3000" {
		t.Errorf("Did not receive expected hostname")
	}
}
