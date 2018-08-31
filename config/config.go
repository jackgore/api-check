package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Hostname string `json:"hostname"`
}

const (
	DefaultConfigFile = ".ac.json"
)

// New creats a new config object from the given filename.
func New(filename string) (Config, error) {
	var conf Config

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		// Dont return an error because its not needed to have a config file
		return conf, nil
	}

	if err := json.Unmarshal(contents, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}
