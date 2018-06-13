package loader

import (
	"os"
	"path/filepath"
)

const (
	extension = ".ac.json"
)

// FindRunDefs finds all run definitions in the working directory
// or in any directory below it.
func FindRunDefs() ([]string, error) {
	files := []string{}

	wd, err := os.Getwd()
	if err != nil {
		return files, err
	}

	filepath.Walk(wd, func(path string, f os.FileInfo, err error) error {
		if err != nil || filepath.Base(path) == extension {
			return nil
		}

		// We are looking to detect the .ac.json file extension

		// Gets the outer .json file extension
		outer := filepath.Ext(path)

		// Gets the string without the .json extension
		inner := filepath.Ext(path[:(len(path) - len(outer))])

		if !f.IsDir() && (inner+outer) == extension {
			files = append(files, path)
		}
		return nil
	})

	return files, nil
}
