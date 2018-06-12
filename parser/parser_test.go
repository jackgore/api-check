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

func TestValidateStatusCode(t *testing.T) {
	// 0 Status code should result in default statuscode being returned
	if code, err := validateStatusCode(0); err != nil {
		t.Errorf("Received unexpected error when validating statuscode")
	} else if code != DefaultStatusCode {
		t.Errorf("Did not receive expected code. Expected %v. Actual:%v", DefaultStatusCode, code)
	}

	// Valid status code should produce itself
	if code, err := validateStatusCode(400); err != nil {
		t.Errorf("Received unexpected error when validating statuscode")
	} else if code != 400 {
		t.Errorf("Did not receive expected code. Expected %v. Actual:%v", 400, code)
	}

	// Invalid status code should throw an error
	if _, err := validateStatusCode(999); err == nil {
		t.Errorf("Expected to receive an error for an invalid statuscode")
	}

	if _, err := validateStatusCode(1); err == nil {
		t.Errorf("Expected to receive an error for an invalid statuscode")
	}

	if _, err := validateStatusCode(-1); err == nil {
		t.Errorf("Expected to receive an error for an invalid statuscode")
	}
}
