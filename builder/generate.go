package builder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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
// Returns the name of the file written.
func CreateSkeletonFile(prefix string) (string, error) {
	if len(prefix) == 0 {
		return "", errors.New("name cannot be empty")
	}

	filename := prefix + extension
	if _, err := os.Stat(filename); err == nil {
		// If the file already exists we must fail to avoid overwriting anything.
		return "", fmt.Errorf("cannot create template file as file %v already exists", filename)
	}

	contents, err := JSONSkeleton()
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(filename, contents, 0644); err != nil {
		return "", err
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

// JSONSkeleteon returns the byte array representation of an empty
// API test definition as part of an array.
func JSONSkeleton() ([]byte, error) {
	skeleton := []APITest{Skeleton()}
	d, err := json.MarshalIndent(&skeleton, "", "    ")
	if err != nil {
		return nil, err
	}

	return d, nil
}
