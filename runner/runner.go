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

	// Ensure the bodies are the same
	if expected.Body != string(body) {
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

func RunTest(test builder.APITest) RunReport {
	report := RunReport{
		Successful: false,
		Test:       test,
	}

	client := &http.Client{} // TODO: Will eventually load a bunch of config from conf file

	u, err := buildURL(test.Hostname, test.Endpoint, test.Request.QueryParams) // TODO: Hostname needs to be allowed to be overwritten in conf file
	if err != nil {
		report.Error = err
		return report
	}

	req, err := http.NewRequest(test.Method, u, bytes.NewBuffer([]byte(test.Request.Body)))
	if err != nil {
		report.Error = err
		return report
	}

	resp, err := client.Do(req)
	if err != nil {
		report.Error = err
		return report
	}

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
