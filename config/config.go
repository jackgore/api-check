package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Hostname string `json:"hostname"`
}

func New(filename string) (Config, error) {
	var conf Config

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return conf, err
	}

	if err := json.Unmarshal(contents, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}
