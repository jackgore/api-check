package runner

import "testing"

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
