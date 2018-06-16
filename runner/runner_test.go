package runner

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/JonathonGore/api-check/builder"
)

const (
	json1 = `{ "testing": "jack" }`
	json2 = `{ "testing": "notjack" }`
	json3 = `[]`
	json4 = `["jack", "hello"]`
	json5 = `["hello", "jack"]`
)

var (
	expected = builder.APIResponse{
		Body: "test",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: http.StatusOK,
	}
)

func TestAssertJSON(t *testing.T) {
	var expected interface{}
	var actual interface{}

	// Both interfaces nil should result in assertion passing
	if !assertJSON(actual, expected) {
		t.Errorf("Received unexpected error when asserting json of nil maps")
	}

	// Only expected nil should pass
	if !assertJSON(actual, expected) {
		t.Errorf("Received unexpected error when asserting json of nil maps")
	}

	if err := json.Unmarshal([]byte(json1), &expected); err != nil {
		t.Errorf("Error unmarshaling while testing: %v", err)
	}

	// The exact same data should result in no error
	if !assertJSON(expected, expected) {
		t.Errorf("Received unexpected error when asserting json of matching maps")
	}

	if err := json.Unmarshal([]byte(json2), &actual); err != nil {
		t.Errorf("Error unmarshaling while testing: %v", err)
	}

	// Mismatching values should fail
	if assertJSON(actual, expected) {
		t.Errorf("Expected to receive error when comparing mismatching maps")
	}

	if err := json.Unmarshal([]byte(json3), &expected); err != nil {
		t.Errorf("Error unmarshaling while testing: %v", err)
	}

	// Two empty arrays should be equal
	if !assertJSON(expected, expected) {
		t.Errorf("Received unexpected error when asserting json of empty arrays")
	}

	if err := json.Unmarshal([]byte(json4), &actual); err != nil {
		t.Errorf("Error unmarshaling while testing: %v", err)
	}

	// Two unequal arrays should fail
	if assertJSON(actual, expected) {
		t.Errorf("Expected to receive error when comparing unequal arrays")
	}

	if err := json.Unmarshal([]byte(json4), &expected); err != nil {
		t.Errorf("Error unmarshaling while testing: %v", err)
	}

	// Order shouldn matter in arrays
	if !assertJSON(actual, expected) {
		t.Errorf("Expected to receive error when comparing reordered arrays")
	}
}

func TestAssertResponse(t *testing.T) {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	validResp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     header,
		Body:       ioutil.NopCloser(bytes.NewBufferString("test")),
	}

	// Valid response should yield successful result
	if success, err := assertResponse(validResp, expected); err != nil {
		t.Errorf("Received unexpected error when asserting response: %v", err)
	} else if !success {
		t.Errorf("Assert response unexpectedly failed")
	}

	invalidCodeResp := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Header:     header,
		Body:       ioutil.NopCloser(bytes.NewBufferString("test")),
	}

	// Mismatching headers in response should yield unsuccessful result
	if success, err := assertResponse(invalidCodeResp, expected); err == nil || success {
		t.Errorf("Expected error and failure when asserting response with invalid code")
	}

	invalidHeaderResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString("test")),
	}

	// Mismatching headers should yield unsuccessful result
	if success, err := assertResponse(invalidHeaderResp, expected); err == nil || success {
		t.Errorf("Expected error and failure when asserting response with mismatching header")
	}

	invalidBodyResp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     header,
		Body:       ioutil.NopCloser(bytes.NewBufferString("mismatching")),
	}

	// Mismatching bodies should yield unsuccessful result
	if success, err := assertResponse(invalidBodyResp, expected); err == nil || success {
		t.Errorf("Expected error and failure when asserting response with mismatching body")
	}
}

func TestBuildQueryString(t *testing.T) {
	var m map[string]string

	// Nil map should produce empty string
	if q := buildQueryString(m); q != "" {
		t.Errorf("Passing empty map should produce empty query. Received: %v", q)
	}

	// Empty map should produce empty string
	m = make(map[string]string)
	if q := buildQueryString(m); q != "" {
		t.Errorf("Passing empty map should produce empty query. Received: %v", q)
	}

	m["key"] = "value"
	if q := buildQueryString(m); q != "?key=value" {
		t.Errorf("Passing map with single value should produce built query. Received: %v", q)
	}

	m["another"] = "key"
	q := buildQueryString(m)
	if q != "?key=value&another=key" && q != "?another=key&key=value" {
		t.Errorf("Passing map with single value should produce built query. Received: %v", q)
	}
}
