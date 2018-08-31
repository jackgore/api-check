package loader

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// the extension used for test definition files.
	extension = ".ac.json"
)

// Consumes a filepath and determines if it has a file extension containing
// two dots (i.e '.') and determines if it matches the given extension.
// For the result to be meaningful ext must contain two '.'.
func hasDoubleDotExt(path, ext string) bool {
	// Gets the outer .json file extension
	outer := filepath.Ext(path)

	// Gets the extension of the string without the .json extension
	inner := filepath.Ext(path[:(len(path) - len(outer))])

	return (inner + outer) == extension
}

// FindRunDefs finds all run definition files in the working directory
// or in any directory below it.
func FindTestDefinitions(dir string) ([]string, error) {
	var files []string

	fmt.Printf("Finding testing definitions for dir %v\n", dir)
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		fmt.Printf("Searching file/dir: %v\n", path)
		if err != nil || filepath.Base(path) == extension {
			// When filepath.Base(path) == extension we are looking at the
			// global configuration file '.ac.json'.
			// Note: if err != nil we currently will not return an error as this
			// will cause all other directories/files to be skipped.
			fmt.Printf("received error when reading file %v: %v\n", path, err)
			return nil
		}

		// We are looking to detect the .ac.json file extension
		if !f.IsDir() && hasDoubleDotExt(path, extension) {
			fmt.Println("found valid file")
			files = append(files, path)
		}
		return nil
	})

	return files, nil
}
