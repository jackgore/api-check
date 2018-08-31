package loader

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var doubleDotExtTests = []struct {
	path      string
	extension string
	result    bool
}{
	{"users.ac.json", ".ac.json", true},
	{"users.c.json", ".ac.json", false},
	{"users.ac.json.json", ".ac.json", false},
	{"users.ac.ac.json", ".ac.json", true},
	{"", ".ac.json", false},
}

func TestHasDoubleDotExt(t *testing.T) {
	for _, test := range doubleDotExtTests {
		if result := hasDoubleDotExt(test.path, test.extension); result != test.result {
			t.Errorf("path %v expected %v for extension %v but received %v",
				test.path, test.result, test.extension, result)
		}
	}
}

func contains(arr []string, str string) bool {
	for _, val := range arr {
		if val == str {
			return true
		}
	}
	return false
}

func TestFindTestDefinitions(t *testing.T) {
	// In order to test this function we create some temporary files in a
	// temporary directory and ensure we can find them once they are created.
	dir, err := ioutil.TempDir("/tmp", "testing")
	if err != nil {
		t.Fatalf("unable to create temporary directory for testing")
	}
	fmt.Printf("Created the following directory: %v", dir)
	defer os.RemoveAll(dir) // clean up all directories used for testing

	// Create a sub directory for testing
	subdir, err := ioutil.TempDir(dir, "testing")
	if err != nil {
		t.Fatalf("unable to create temporary directory for testing")
	}

	tmpfile, err := ioutil.TempFile(dir, "*.ac.json")
	if err != nil {
		t.Fatalf("unable to create temporary file for testing")
	}

	subtmpfile, err := ioutil.TempFile(subdir, "*.ac.json")
	if err != nil {
		t.Fatalf("unable to create temporary file for testing")
	}

	_, err = ioutil.TempFile(subdir, "testing")
	if err != nil {
		t.Fatalf("unable to create temporary file for testing")
	}

	files, err := FindTestDefinitions(dir)
	if err != nil {
		t.Fatalf("received error %v when trying to find test definition", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected to only find two test definition files - found %v", len(files))
	}

	if !contains(files, tmpfile.Name()) || !contains(files, subtmpfile.Name()) {
		t.Fatalf("did not receive expected files from FindTestDefinitions")
	}
}
