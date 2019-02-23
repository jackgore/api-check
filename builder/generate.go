package builder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	extension            = ".ac.json" // TODO: extension should be globalized.
	SkeletonDescription  = "Test Description"
	SkeletonEndpoint     = "/"
	SkeletonHostname     = "http://localhost"
	SkeletonMethod       = "GET"
	SkeletonResponseCode = 200
)

// CreateSkeletonFile writes a .ac.json skeleton file with the given prefix.
// Returns the name of the file written and an error if applicable.
func CreateSkeletonFile(prefix string) (string, error) {
	if len(prefix) == 0 {
		return "", errors.New("filename cannot be empty")
	}

	filename := prefix

	// If the provided prefix contains the '.ac.json' suffic don't append
	// an addition '.ac.json' suffix.
	if !strings.HasSuffix(prefix, extension) {
		filename = filename + extension
	}

	// If the file already exists we must fail to avoid overwriting anything.
	if _, err := os.Stat(filename); err == nil {
		return "", fmt.Errorf("cannot create template file as file %v already exists", filename)
	}

	contents, err := JSONSkeleton()
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(filename, contents, 0644); err != nil {
		return "", fmt.Errorf("unable to write skeleton file %v: %v", filename, err)
	}

	return filename, nil
}

// Skeleton creates an APITest skeleton to use for generating a skeleton file.
func Skeleton() APITest {
	// TODO: Might be cool to read in `.ac.json` and use those values for the skeleton.
	return APITest{
		Description: SkeletonDescription,
		Endpoint:    SkeletonEndpoint,
		Hostname:    SkeletonHostname,
		Method:      SkeletonMethod,
		Request: APIRequest{
			Headers:     make(map[string]string),
			QueryParams: make(map[string]string),
		},
		Response: APIResponse{
			Headers:    make(map[string]string),
			StatusCode: SkeletonResponseCode,
		},
	}
}

// JSONSkeleteon returns the byte array representation of an empty API test
// definition as part of an array.
func JSONSkeleton() ([]byte, error) {
	skeleton := []APITest{Skeleton()}

	// Marshal the into JSON skeleton using nice indentation.
	d, err := json.MarshalIndent(&skeleton, "", "    ")
	if err != nil {
		return nil, err
	}

	return d, nil
}
