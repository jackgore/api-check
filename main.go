package main

import (
	"github.com/JonathonGore/api-check/suite"
)

const (
	confFile         = ".ac.json"
	defaultVerbosity = true
)

func main() {
	suite.Verbose(defaultVerbosity)
	suite.RunStandalone()
}
