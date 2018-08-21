package runner

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/JonathonGore/api-check/builder"
)

const (
	json1 = `{ "testing": "jack" }`
	json2 = `{ "testing": "notjack" }`
	json3 = `[]`
	json4 = `["jack", "hello"]`
	json5 = `["hello", "jack"]`
	json6 = `{ "testing": "jack", "extra":"key" }`
)

var (
	basicAPI = builder.APIResponse{
		Body: "test",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: http.StatusOK,
	}

	noBodyAPI = builder.APIResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: http.StatusOK,
	}

	emptyJSONAPI = builder.APIResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		JSON:       "",
		StatusCode: http.StatusOK,
	}

	basicResponse = http.Response{
		StatusCode: http.StatusOK,
		// Header set in init() function
		Body: ioutil.NopCloser(bytes.NewBufferString("test")),
	}

	noBodyResponse = http.Response{
		StatusCode: http.StatusOK,
		// Header set in init() function
		Body: ioutil.NopCloser(bytes.NewBufferString("")),
	}

	statusCodeResponse = http.Response{
		StatusCode: http.StatusUnauthorized,
		// Header set in init() function
		Body: ioutil.NopCloser(bytes.NewBufferString("test")),
	}

	headerResponse = http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString("test")),
	}

	bodyResponse = http.Response{
		StatusCode: http.StatusOK,
		// Header set in init() function
		Body: ioutil.NopCloser(bytes.NewBufferString("mismatching")),
	}
)

func init() {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	basicResponse.Header = header
	statusCodeResponse.Header = header
	bodyResponse.Header = header
	noBodyResponse.Header = header
}

var buildRequestTests = []struct {
	test builder.APITest
	expected *http.Request
	success bool
}{

}

func TestBuildRequest(t *testing.T) {
	for _, test := range buildRequestTests {
		r, err := buildRequest(test.test)
		if test.success && err != nil {
			t.Errorf("expected test to succeed but it failed with error: %v", err)
		}

		if !test.success && err == nil {
			t.Errorf("expected test to fail but it succeeded: %v", err)
		}

		if !reflect.DeepEqual(*r, *test.expected) {
			t.Errorf("mismatching request");
		}
	}
}

var assertJSONTests = []struct {
	actual   string
	expected string
	succeed  bool
}{
	{"{}", "{}", true},    // Two empty responses should be equal
	{json1, json1, true},  // Equal JSON should succeed
	{json1, json2, false}, // Mismatching values for a JSON key should fail
	{json3, json3, true},  // Empty arrays should be equal
	{json3, json4, false}, // If we dont get the expected array elements we should fail
	{json5, json4, false}, // Bad ordering of arrays should cause a failure
	{json1, json6, false}, // Extra key in expected should fail
	{json6, json1, true},  // Extra key in actual should succeed
}

func TestAssertJSON(t *testing.T) {
	var actual interface{}
	var expected interface{}

	for _, test := range assertJSONTests {
		if err := json.Unmarshal([]byte(test.actual), &actual); err != nil {
			t.Errorf("Unable to unmarshal JSON: %v", err)
		}

		if err := json.Unmarshal([]byte(test.expected), &expected); err != nil {
			t.Errorf("Unable to unmarshal JSON: %v", err)
		}

		if assertJSON(actual, expected) != test.succeed {
			succeedText := "passed"
			if !test.succeed {
				succeedText = "failed"
			}
			t.Errorf("Received: %v Expected: %v - Test should have: %v", test.actual, test.expected, succeedText)
		}
	}
}

var assertResponseTests = []struct {
	actual   *http.Response
	expected builder.APIResponse
	succeed  bool
}{
	{&basicResponse, basicAPI, true},       // Matching response and expected should succeed
	{&statusCodeResponse, basicAPI, false}, // Mismatching status codes should fail
	{&headerResponse, basicAPI, false},     // Missing header should fail
	{&bodyResponse, basicAPI, false},       // Mismatching body should fail
	{&noBodyResponse, noBodyAPI, true},     // Actual and expected with no body's should succeed
	{&bodyResponse, noBodyAPI, true},       // If we dont expect a body but still receive one then succeed
}


func TestAssertResponse(t *testing.T) {
	for _, test := range assertResponseTests {
		if ok, _ := assertResponse(test.actual, test.expected); ok != test.succeed {
			succeedText := "passed"
			if !test.succeed {
				succeedText = "failed"
			}
			t.Errorf("Received: %v Expected: %+v - Test should have: %+v", test.actual, test.expected, succeedText)
		}
	}
}

var buildQueryStringTests = []struct {
	input    map[string]string
	expected []string
}{
	{map[string]string{}, []string{""}},
	{map[string]string{"key": "value"}, []string{"?key=value"}},
	{map[string]string{"key": "value", "another": "key"}, []string{"?another=key&key=value", "?key=value&another=key"}},
}

func TestBuildQueryString(t *testing.T) {
	for _, test := range buildQueryStringTests {
		q := buildQueryString(test.input)
		passed := false

		for _, expected := range test.expected {
			if expected == q {
				passed = true
				break
			}
		}

		if !passed {
			t.Errorf("Received: %v. Expected something similar to: %v", q, test.expected[0])
		}
	}
}
