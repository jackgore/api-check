package runner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/JonathonGore/api-check/builder"
)

type RunReport struct {
	Test           builder.APITest
	Successful     bool
	Error          error
	FailureMessage string
}

// Consumes a map of string => string representing query params
// and builds the query string in the form "?<key>=<value>&<key>=<value>"
func buildQueryString(query map[string]string) string {
	if len(query) == 0 {
		return ""
	}

	qstring := "?"
	for key, val := range query {
		if qstring != "?" {
			qstring = qstring + "&"
		}

		qstring = qstring + key + "=" + val
	}

	return qstring
}

// TODO: This should maybe return http.URL
func buildURL(hostname, endpoint string, query map[string]string) (string, error) {
	qstring := buildQueryString(query)

	return hostname + endpoint + qstring, nil
}

// AssertResponse consume the http response from the server and the struct containing the
// expected results and compares the two and ensures they are equal
func assertResponse(resp *http.Response, expected builder.APIResponse) (bool, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// Ensure status code is what is expected
	if expected.StatusCode != resp.StatusCode {
		return false, fmt.Errorf("Unexpected status code received\n\nExpected:\n%v\n\nActual:\n%v\n\n", expected.StatusCode, resp.StatusCode)
	}

	// Ensure the bodies are the same only if the expected body is non-empty
	// NOTE: Right now we have no way of asserting the response body is empty
	if expected.Body != "" && expected.Body != string(body) {
		return false, fmt.Errorf("Mismatching bodies\n\nExpected:\n%v\n\nActual:\n%v\n\n", expected.Body, string(body))
	}

	// Ensure headers are what we expect
	for key, value := range expected.Headers {
		if value != resp.Header.Get(key) {
			return false, fmt.Errorf("Mismatching %v header\n\nExpected:\n%v\n\nActual:\n%v\n\n", key, value, resp.Header.Get(key))
		}
	}

	return true, nil
}

// BuildRequest consumes an api test object and produces the corresponding http request
// that will be sent by the http client
func buildRequest(test builder.APITest) (*http.Request, error) {
	u, err := buildURL(test.Hostname, test.Endpoint, test.Request.QueryParams)
	if err != nil {
		return nil, err
	}

	// Build request object attaching the specified method, url and body
	req, err := http.NewRequest(test.Method, u, bytes.NewBuffer([]byte(test.Request.Body)))
	if err != nil {
		return nil, err
	}

	// Attach the specified request headers
	for key, value := range test.Request.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

// RunTest consumes a test object and runs the test against the configured server
// produces a RunReport of the results
func RunTest(test builder.APITest) RunReport {
	report := RunReport{
		Successful: false,
		Test:       test,
	}

	client := &http.Client{} // TODO: Will eventually load a bunch of config from conf file

	req, err := buildRequest(test)
	if err != nil {
		report.Error = err
		return report
	}

	resp, err := client.Do(req)
	if err != nil {
		report.Error = err
		return report
	}

	// Compare result to expected result
	report.Successful, report.Error = assertResponse(resp, test.Response)

	return report
}

func RunTests(tests []builder.APITest) []RunReport {
	reports := make([]RunReport, len(tests))

	for i, test := range tests {
		reports[i] = RunTest(test)
	}

	return reports
}
