package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/JonathonGore/api-check/builder"
)

func Parse(filename string) ([]builder.APITest, error) {
	tests := []builder.APITest{}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return tests, err
	}

	err = json.Unmarshal(contents, &tests)
	if err != nil {
		return tests, err
	}

	for i, test := range tests {
		if err = validate(test); err != nil {
			return tests, fmt.Errorf("Error in test #%v: %v", i+1, err)
		}
	}

	return tests, nil
}

func validate(test builder.APITest) error {
	return nil
}
