package parser

import (
	"net/http"
	"testing"

	"github.com/JonathonGore/api-check/config"
)

var p = New(config.Config{})

func TestValidateMethod(t *testing.T) {
	// Valid http method should be validate
	if method, err := p.validateMethod("GET"); err != nil {
		t.Errorf("Received unexpected error when validating valid method")
	} else if method != http.MethodGet {
		t.Errorf("Did not receive expected method")
	}

	// Valid http method should be validate but be uppercase
	if method, err := p.validateMethod("get"); err != nil {
		t.Errorf("Received unexpected error when validating valid method")
	} else if method != http.MethodGet {
		t.Errorf("Did not receive expected method")
	}

	if method, err := p.validateMethod("oPtiONs"); err != nil {
		t.Errorf("Received unexpected error when validating valid method")
	} else if method != http.MethodOptions {
		t.Errorf("Did not receive expected method")
	}

	// Invalid method should fail
	if _, err := p.validateMethod("fake method"); err == nil {
		t.Errorf("Expected to receive error when validating invalid http method")
	}

	// Empty method should be set to default value
	if method, err := p.validateMethod(""); err != nil {
		t.Errorf("Received unexpected error when validating valid method")
	} else if method != http.MethodGet {
		t.Errorf("Did not receive expected method")
	}
}

func TestValidateEndpoint(t *testing.T) {
	// Empty endpoint should be set to default value
	if endpoint, err := p.validateEndpoint(""); err != nil {
		t.Errorf("Received unexpected error when validating empty endpoint")
	} else if endpoint != DefaultEndpoint {
		t.Errorf("Empty endpoint should produce default endpoint")
	}

	// Endpoint without leading slash should fail
	if _, err := p.validateEndpoint("testing"); err == nil {
		t.Errorf("Expected to receive error when validating endpoint without leading slash")
	}

	// Valid endpoint should produce itself
	if endpoint, err := p.validateEndpoint("/users/apicheck"); err != nil {
		t.Errorf("Received unexpected error when validating endpoint")
	} else if endpoint != "/users/apicheck" {
		t.Errorf("Did not receive expected endpoint")
	}
}

func TestValidateHostname(t *testing.T) {
	hostname := "http://localhost:3000"

	// Empty hostname should be an error
	if _, err := p.validateHostname(""); err == nil {
		t.Errorf("Expected to receive error when validating empty hostname")
	}

	// Creating a parser with a config with hostname should override whatever is passed in
	configParser := Parser{conf: config.Config{Hostname: hostname}}
	if result, err := configParser.validateHostname(""); err != nil {
		t.Errorf("Received unexpected error when validating hostname")
	} else if result != hostname {
		t.Errorf("Did not receive expected hostname")
	}

	if result, err := configParser.validateHostname("https://www.google.com"); err != nil {
		t.Errorf("Received unexpected error when validating hostname")
	} else if result != hostname {
		t.Errorf("Did not receive expected hostname")
	}

	// Valid hostname should succeed
	if result, err := p.validateHostname(hostname); err != nil {
		t.Errorf("Received unexpected error when validating hostname")
	} else if result != hostname {
		t.Errorf("Did not receive expected hostname")
	}
}

func TestValidateStatusCode(t *testing.T) {
	// 0 Status code should result in default statuscode being returned
	if code, err := p.validateStatusCode(0); err != nil {
		t.Errorf("Received unexpected error when validating statuscode")
	} else if code != DefaultStatusCode {
		t.Errorf("Did not receive expected code. Expected %v. Actual:%v", DefaultStatusCode, code)
	}

	// Valid status code should produce itself
	if code, err := p.validateStatusCode(400); err != nil {
		t.Errorf("Received unexpected error when validating statuscode")
	} else if code != 400 {
		t.Errorf("Did not receive expected code. Expected %v. Actual:%v", 400, code)
	}

	// Invalid status code should throw an error
	if _, err := p.validateStatusCode(999); err == nil {
		t.Errorf("Expected to receive an error for an invalid statuscode")
	}

	if _, err := p.validateStatusCode(1); err == nil {
		t.Errorf("Expected to receive an error for an invalid statuscode")
	}

	if _, err := p.validateStatusCode(-1); err == nil {
		t.Errorf("Expected to receive an error for an invalid statuscode")
	}
}
