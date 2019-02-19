package builder

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var createSkeletonFileTests = []struct {
	prefix string
	err    error
}{
	{"", errors.New("name cannot be empty")}, // Empty prefix should cause a failure to occur.
	{"unit-test.ac.json", nil},               // Empty prefix should cause a failure to occur.
	{"unit-test", nil},                       // Empty prefix should cause a failure to occur.
}

func TestCreateSkeletonFile(t *testing.T) {
	for _, test := range createSkeletonFileTests {
		err := CreateSkeletonFile(test.prefix)

		if test.err != nil {
			// Case where we expect an error from the test.
			if err == nil {
				t.Errorf("Expected test to fail with error: %v but got nil", test.err)
			} else if err.Error() != test.err.Error() {
				t.Errorf("Expected test to fail with error: %v but got: %v", test.err, err)
			}
		} else if test.err == nil {
			// Case where we expect the test to succeed.
			if err != nil {
				t.Errorf("Expected test to succeed but failed with error: %v", err)
			} else {
				// Make sure the contents of the file match those expected
				contents, err := ioutil.ReadFile("unit-test.ac.json")
				if err != nil {
					t.Errorf("Unable to complete test: %v", err)
				}
				expectedContents, _ := JSONSkeleton()
				if !reflect.DeepEqual(contents, expectedContents) {
					t.Errorf("Found mismatching results when comparing skeleton contents")
				}
			}
		}

		os.Remove("unit-test.ac.json")
	}
}
