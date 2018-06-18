package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/JonathonGore/api-check/builder"
	"github.com/JonathonGore/api-check/config"
)

const (
	DefaultEndpoint   = "/"
	DefaultMethod     = http.MethodGet
	DefaultStatusCode = http.StatusOK
)

type Parser struct {
	conf    config.Config
	methods map[string]bool
}

// New consumes a api-check config object and builds a new
// Parser object.
func New(conf config.Config) Parser {
	methods := map[string]bool{
		http.MethodGet:     true,
		http.MethodHead:    true,
		http.MethodPost:    true,
		http.MethodPut:     true,
		http.MethodPatch:   true,
		http.MethodDelete:  true,
		http.MethodConnect: true,
		http.MethodOptions: true,
		http.MethodTrace:   true,
	}

	return Parser{
		conf:    conf,
		methods: methods,
	}
}

// ParseFile consumes a single filename and parses it into a list
// of api tests.
func (p *Parser) ParseFile(file string) ([]builder.APITest, error) {
	tests := []builder.APITest{}

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return tests, err
	}

	// TODO: It would be cool if we could detect extra fields and warn the user about them
	err = json.Unmarshal(contents, &tests)
	if err != nil {
		return tests, err
	}

	for i, test := range tests {
		if tests[i], err = p.validate(test); err != nil {
			return tests, fmt.Errorf("error in test #%v: %v", i+1, err)
		}
	}

	return tests, nil
}

// Parse consumes a list of filenames and attempts to parse them into a list
// of api tests.
func (p *Parser) Parse(filenames []string) ([]builder.APITest, error) {
	tests := []builder.APITest{}

	for _, name := range filenames {
		results, err := p.ParseFile(name)
		if err != nil {
			return tests, fmt.Errorf("Error parsing file: %v - %v", name, err)
		}

		tests = append(tests, results...)
	}

	return tests, nil
}

// Validates the given status code. Defaulting to the default value
// if needed.
func (p *Parser) validateStatusCode(code int) (int, error) {
	// If code is 0 this means it was not set - so reset to default
	if code == 0 {
		return DefaultStatusCode, nil
	}

	if code < 100 || code >= 600 {
		return code, fmt.Errorf("HTTP status code out of range")
	}

	return code, nil
}

// ValidMethod consumes a string and asserts that it is a valid http method
// if it is not an error is returned. If the given method is empty it returns
// the default method. Case insensitive.
func (p *Parser) validateMethod(method string) (string, error) {
	if len(method) == 0 {
		return DefaultMethod, nil
	}

	method = strings.ToUpper(method) // Uppercase to avoid case issues

	// We need method to be a http supported method
	if _, ok := p.methods[method]; !ok {
		return method, fmt.Errorf("Received unsupport http method: %v", method)
	}

	return method, nil
}

// Validates the given hostname and assert it is either non-empty or specified
// in the api-check config.
// Hostname is a required field for api-check
// Whatever is in the test object should override the conf
func (p *Parser) validateHostname(hostname string) (string, error) {
	if len(hostname) == 0 {
		if p.conf.Hostname != "" {
			return p.conf.Hostname, nil
		}

		return "", fmt.Errorf("hostname is a required field")
	}

	u, err := url.ParseRequestURI(hostname)
	if err != nil {
		return "", fmt.Errorf("malformed hostname provided")
	}

	// TODO: not sure how much validation we want to do - this currently will allow schemes
	// that are not http or https

	// Reformat the url to ensure there is not extra text like a trailing slash.
	// NOTE: the prevents people from doing things like http://localhost:3000/v1
	// Not sure if we want to allow that or not.
	return fmt.Sprintf("%v://%v", u.Scheme, u.Host), nil
}

// Validates the given endpoint returning either the input string,
// a default value should the input be empty or an error if input is invalid.
func (p *Parser) validateEndpoint(endpoint string) (string, error) {
	// We need the endpoint to begin with a forward slasg '/'
	// Note: Not sure if we want to add a leading slach if the slash is missing
	if len(endpoint) == 0 {
		return DefaultEndpoint, nil
	}

	if string(endpoint[0]) != "/" {
		return endpoint, fmt.Errorf("Provided endpoint must begin witha '/'. Found: %v.", endpoint[0])
	}

	return endpoint, nil
}

// Validate is used to validate paramaters of an APItest and replace empty
// paramaters with default/initialized values.
func (p *Parser) validate(test builder.APITest) (builder.APITest, error) {
	var err error

	test.Endpoint, err = p.validateEndpoint(test.Endpoint)
	if err != nil {
		return test, err
	}

	test.Hostname, err = p.validateHostname(test.Hostname)
	if err != nil {
		return test, err
	}

	test.Response.StatusCode, err = p.validateStatusCode(test.Response.StatusCode)
	if err != nil {
		return test, err
	}

	test.Method, err = p.validateMethod(test.Method)
	if err != nil {
		return test, err
	}

	return test, nil
}
