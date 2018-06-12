package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JonathonGore/api-check/builder"
)

const (
	DefaultEndpoint = "/"
	DefaultStatusCode = http.StatusOK
)

func Parse(filename string) ([]builder.APITest, error) {
	tests := []builder.APITest{}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return tests, err
	}

	err = json.Unmarshal(contents, &tests)
	if err != nil {
		return tests, err
	}

	for i, test := range tests {
		if tests[i], err = validate(test); err != nil {
			return tests, fmt.Errorf("Error in test #%v: %v", i+1, err)
		}
	}

	return tests, nil
}

// Validates the given status code. Defaulting to the default value
// if needed.
func validateStatusCode(code int) (int, error) {
	// If code is 0 this means it was not set - so reset to default
	if code == 0 {
		return DefaultStatusCode, nil
	}

	if code < 100 || code >= 600 {
		return code, fmt.Errorf("HTTP status code out of range")
	}

	return code, nil
}

// Validates the given hostname and assert it is either non-empty or specified
// in the api-check config.
// Hostname is a required field for api-check
func validateHostname(hostname string) (string, error) {
	if len(hostname) == 0 {
		// TODO: We need to check the config file for a hostname once we enable that functionality
		return "", fmt.Errorf("Hostname is a required field in order to run api-check")
	}

	// TODO we need to validate hostname is in valid format - will use net/http/url
	return hostname, nil
}

// Validates the given endpoint returning either the input string,
// a default value should the input be empty or an error if input is invalid.
func validateEndpoint(endpoint string) (string, error) {
	// We need the endpoint to begin with a forward slasg '/'
	// Note: Not sure if we want to add a leading slach if the slash is missing 
	if len(endpoint) == 0 {
		return DefaultEndpoint, nil
	}

	// TODO: Probably want to use the 'filepath' package to do this
	if string(endpoint[0]) != "/" {
		return endpoint, fmt.Errorf("Provided endpoint must begin witha '/'. Found: %v.", endpoint[0])
	}

	return endpoint, nil
}

// Validate is used to validate paramaters of an APItest and replace empty
// paramaters with default/initialized values.
func validate(test builder.APITest) (builder.APITest, error) {
	var err error

	test.Endpoint, err = validateEndpoint(test.Endpoint)
	if err != nil {
		return test, err
	}

	test.Hostname, err = validateHostname(test.Hostname)
	if err != nil {
		return test, err
	}

	test.Response.StatusCode, err = validateStatusCode(test.Response.StatusCode)
	if err != nil {
		return test, err
	}

	return test, nil
}
