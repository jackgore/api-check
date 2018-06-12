package runner

import (
	"bytes"
	"net/http"

	"github.com/JonathonGore/api-check/builder"
)

type RunReport struct {
	Successes int
	Failures int
	Errors int
}

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

func assertResponse(resp *http.Response, expected builder.APIResponse) (bool, error) {
	return true, nil
}

func RunTest(test builder.APITest) (bool, error) {
	client := &http.Client{} // TODO: Will eventually load a bunch of config from conf file

	u, err := buildURL(test.Hostname, test.Endpoint, test.Request.QueryParams) // TODO: Hostname needs to be allowed to be overwritten in conf file
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(test.Method, u, bytes.NewBuffer([]byte(test.Request.Body)))
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	return assertResponse(resp, test.Response)
}

func RunTests(tests []builder.APITest) RunReport {
	report := RunReport{}

	for _, test := range tests {
		success, err := RunTest(test)
		if err != nil {
			report.Errors = report.Errors + 1
		} else if success {
			report.Successes = report.Successes + 1
		} else {
			report.Failures = report.Failures + 1
		}
	}

	return report
}
