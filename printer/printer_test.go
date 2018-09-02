package printer

import (
	"testing"

	"github.com/JonathonGore/api-check/builder"
)

func TestSucceededText(t *testing.T) {
	if result := succeededText(true); result != "succeeded" {
		t.Errorf("expected succeeded but received: %v", result)
	}

	if result := succeededText(false); result != "failed" {
		t.Errorf("expected failed but received: %v", result)
	}
}

var buildDescriptionTests = []struct {
	test   builder.APITest
	result string
}{
	{builder.APITest{Description: "test"}, "test"},
	{builder.APITest{Hostname: "testing", Endpoint: "endpoint", Description: "test"}, "test"},
	{builder.APITest{Hostname: "localhost", Endpoint: "/apps"}, "localhost/apps"},
}

func TestBuildDescription(t *testing.T) {
	for _, test := range buildDescriptionTests {
		if result := buildDescription(test.test); result != test.result {
			t.Errorf("expected %v but received %v", test.result, result)
		}
	}
}
