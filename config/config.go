package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config is used to specifiy global config used by the api-check CLI and Go
// test runner. The api-check CLI will look a `.ac.json` file in the same
// directory as where the command is ran to override normal functionality.
type Config struct {
	// Default hostname used for API request, prevents you from needing to
	// specifiy this for each test definition.
	Hostname string `json:"hostname"`

	// Name of a bash script to execute before running the test suite.
	SetupScript string `json:"setup-script"`

	// Name of a bash script to execute before finishing the test suite.
	CleanupScript string `json:"cleanup-script"`

	// MuteScriptOutput determines if the output from setup and cleanup script
	// should be surpressed.
	MuteScriptOutput bool `json:"mute-script-output"`
}

const (
	// The default file api-check will look for config in.
	DefaultConfigFile = ".ac.json"

	// The default option for muting script output.
	DefaultMuteScriptOutput = false
)

// DefaultConfig we will use for the app.
// NOTE: Currently this is not really needed as all values in config default to
// the value we want to use (i.e Golang zero values).
var DefaultConfig = Config{
	MuteScriptOutput: DefaultMuteScriptOutput,
}

// New creats a new config object from the given filename.
func New(filename string) (Config, error) {
	var conf Config

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		// Dont return an error because its not needed to have a config file
		return DefaultConfig, nil
	}

	if err := json.Unmarshal(contents, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}
