package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/JonathonGore/api-check/builder"
)

// RunReport describes the result of a single test.
type RunReport struct {
	Test           builder.APITest
	Successful     bool
	Error          error
	FailureMessage string
}

// buildQueryString Consumes a map of string => string representing query params
// and builds the query string in the form "?<key>=<value>&<key>=<value>".
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

const (
	stringType = "string"
	numberType = "number"
	intType    = "int"
	boolType   = "boolean"
)

// assertJSONType consumes an interface of "unknown" type and asserts that its
// underlying type is that described by the expectedType string.
func assertJSONType(value interface{}, expectedType string) bool {
	expectedType = strings.ToLower(expectedType) // Avoid case issues.

	switch expectedType {
	case stringType:
		// Value needs to be able to be unmarshaled as a string.
		if _, ok := value.(string); !ok {
			return false
		}
	case numberType:
		// If number type is specified it can either be a "double" or an int.
		if _, ok := value.(float64); !ok {
			if _, ok = value.(int); !ok {
				return false
			}
		}
	case intType:
		if _, ok := value.(int); !ok {
            // Sometimes for some reason the underlying type is a float, but is
            // really only an int
            if val, ok := value.(float64); !ok || val != float64(int(val)) {
                return false
             }
		}
	case boolType:
		if _, ok := value.(bool); !ok {
			return false
		}
	default:
		return false // Unexpected type
	}

	return true
}

// assertJSONArray asserts that the given interface is an array of the provided
// object type. 
func assertJSONArray(actual interface{}, expected interface{}) bool {
    // Need to make sure actual is an array
    values, ok := actual.([]interface{})
    if !ok {
        return false
    } else if len(values) == 0 {
        // Empty list is trivially true.
        return true
    }

    // For each value we need to assert its JSONStructure.
    for _, val := range values {
        // Check to see if expected is a map
        if values, ok := expected.(map[string]interface{}); ok {
            if !assertJSONStructure(val, values) {
                return false
            }
        } else if expectedType, ok := expected.(string); ok {
            if !assertJSONType(val, expectedType) {
                return false
            }
        } else {
            // Unexpected type in arrayOf
            return false
        }
     }

     return true
}

// assertJSONStructure consumes the actual response from the server and asserts
// that it has the exact keys as specified in the provided TypeOf.
// TODO: Right now we allow there to be extra keys in the actual
// response it may be useful to allow this to be configurable.
func assertJSONStructure(actual interface{}, expected interface{}) bool {
	if len(expected) == 0 {
		return true
	}

	// Note: as of right now JSONType should not be an array, which means in
	// order for this function to return true actual must be able to be marshaled
	// as a map[string]interface
	actualMap, ok := actual.(map[string]interface{})
	if !ok {
		return false
	}

	// Now for each key in expected we need to make sure it exists in actual
	// and potentially assertJSONStructure on the value.
	for key, val := range expected {
		// If the key is not in actual we have failed the test
		actualVal, ok := actualMap[key]
		if !ok {
			return false
		}

		// If val is a string, it should be describing the expected type
		if expectedType, ok := val.(string); ok {
			// Now we need to make sure actualVal is the correct type.
			if !assertJSONType(actualVal, expectedType) {
				return false
			}
		} else if structure, ok := val.(builder.JSONType); ok {
			// Need to recursively ensure this is true.
			if !assertJSONStructure(actualVal, structure) {
				return false
			}
        } else if arrayType, ok := val.(builder.JSONArrayOf); ok {
           // Expected an array of the given type
           if !assertJSONArray(actualVal, arrayType.Value) {
                return false
           }
		} else {
			// None of the options were matched
			return false
		}
	}
	return true
}

// Asserts that the actual and expected JSON are equal.
// Behaviour is defined such that should there be extra keys in the actual map that is ok,
// so long as every key present in expected is in actual with the same value.
func assertJSON(actual interface{}, expected interface{}) bool {
	if expected == nil {
		return true
	}

	expectedMap, ok := expected.(map[string]interface{})
	if ok {
		// If expected is a map with no keys - pass the test
		if len(expectedMap) == 0 {
			return true
		}

		// If acutal is not a json object return false
		actualMap, ok := actual.(map[string]interface{})
		if !ok {
			return false // TODO: Include error messages
		}

		for key, exp := range expectedMap {
			if acc, ok := actualMap[key]; !ok {
				return false
			} else if !assertJSON(acc, exp) {
				return false
			}
		}

		return true
	}

	// TODO: consider allowing arrays to be in different orders
	return reflect.DeepEqual(actual, expected)
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

	// NOTE: There are basically 3 ways for us to compare request body content.
	// 1) First is by directly comapring the content bytes to what is stored
	//    in the `Body` variable in the test.
	// 2) Second is by directly comparing JSON content, by comparing the JSON
	//    provided with what is provided in the HTTP response.
	// 3) Third is by only looking at the structure of the returned JSON.
	//
	// It is important that we only perform one of these three actions as it
	// could result in weird behaviour if we do otherwise.

	// Ensure the bodies are the same only if the expected body is non-empty
	// NOTE: Right now we have no way of asserting the response body is empty
	if expected.Body != "" && expected.Body != string(body) {
		return false, fmt.Errorf("Mismatching bodies\n\nExpected:\n%v\n\nActual:\n%v\n\n", expected.Body, string(body))
	}

	// Check the structure of the response if TypeOf is present in API Test
	if expected.Body == "" && expected.TypeOf != nil {
		var actual interface{}

		err = json.Unmarshal(body, &actual)
		if err != nil {
			return false, fmt.Errorf("received JSON in unexpected format %v", err)
		}

		if !assertJSONStructure(actual, expected.TypeOf) {
			return false, fmt.Errorf("mismatching JSON structure")
		}
	}

	// Only assert JSON if defined and body is not
	if len(expected.TypeOf) == 0 && expected.Body == "" && expected.JSON != nil {
		var actual interface{}

		err = json.Unmarshal(body, &actual)
		if err != nil {
			return false, fmt.Errorf("Received unexpected error when unmarshaling JSON %v", err)
		}

		if !assertJSON(actual, expected.JSON) {
			return false, fmt.Errorf("Mismatching JSON")
		}
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
// that will be sent by the http client to the server.
func buildRequest(test builder.APITest) (*http.Request, error) {
	u, err := buildURL(test.Hostname, test.Endpoint, test.Request.QueryParams)
	if err != nil {
		return nil, err
	}

	var buffer *bytes.Buffer

	// Only attach json to body if its non-nil with at least 1 key
	if test.Request.JSON != nil {
		contents, err := json.Marshal(test.Request.JSON)
		if err != nil {
			return nil, err
		}
		buffer = bytes.NewBuffer(contents)
	} else {
		buffer = bytes.NewBuffer([]byte(test.Request.Body))
	}

	// Build request object attaching the specified method, url and body
	req, err := http.NewRequest(test.Method, u, buffer)
	if err != nil {
		return nil, err
	}

	// Attach the specified request headers
	for key, value := range test.Request.Headers {
		req.Header.Set(key, value)
	}

	// Determine if we need to set a cookie header
	if test.Request.Cookies != nil && len(test.Request.Cookies) > 0 {
		cookieHeader := ""
		for i, cookie := range test.Request.Cookies {
			if i > 0 {
				cookieHeader = cookieHeader + "; "
			}
			cookieHeader = cookieHeader + fmt.Sprintf("%v=%v", cookie.Name, cookie.Value)
		}
		req.Header.Set("Cookie", cookieHeader)
	}

	return req, nil
}

// RunTest consumes an API test to be run against the configured server
// produces a RunReport of the results of the test.
func RunTest(test builder.APITest) RunReport {
	report := RunReport{
		Successful: false,
		Test:       test,
	}

	// TODO: Will eventually load a bunch of http client config (i.e. custom timeout)
	client := &http.Client{}

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

	report.Successful, report.Error = assertResponse(resp, test.Response)

	return report
}

// RunTests consumes a slice of APITests, runs each test and produces
// a slice of RunReports for each test that is ran.
func RunTests(tests []builder.APITest) []RunReport {
	reports := make([]RunReport, len(tests))

	for i, test := range tests {
		reports[i] = RunTest(test)
	}

	return reports
}
