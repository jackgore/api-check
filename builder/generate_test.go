package builder

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var createSkeletonFileTests = []struct {
	prefix   string
	filename string
	err      error
}{
	{"", "", errors.New("filename cannot be empty")}, // Empty prefix should cause a failure to occur.
	{"unit-test.ac.json", "unit-test.ac.json", nil},  // Empty prefix should cause a failure to occur.
	{"unit-test", "unit-test.ac.json", nil},          // Empty prefix should cause a failure to occur.
}

func TestCreateSkeletonFile(t *testing.T) {
	for _, test := range createSkeletonFileTests {
		filename, err := CreateSkeletonFile(test.prefix)
		if (err != nil && test.err == nil) || (err == nil && test.err != nil) ||
			(err != nil && err.Error() != test.err.Error()) {
			t.Errorf("Expected error: %v but received: %v", test.err, err)
			continue
		}

		if test.filename != filename {
			t.Errorf("Expected filename: %v but received: %v", test.filename, filename)
			continue
		}

		// If we expect the test to succeed check the contents of the file.
		if test.err == nil {
			contents, err := ioutil.ReadFile("unit-test.ac.json")
			if err != nil {
				t.Errorf("Unable to complete test: %v", err)
				continue
			}

			expectedContents, _ := JSONSkeleton()
			if !reflect.DeepEqual(contents, expectedContents) {
				t.Errorf("Found mismatching results when comparing skeleton contents")
			}
		}

		os.Remove(filename)
	}
}
