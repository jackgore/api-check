package suite

import (
	"fmt"
	"os"
	"testing"

	"github.com/JonathonGore/api-check/config"
	"github.com/JonathonGore/api-check/loader"
	"github.com/JonathonGore/api-check/parser"
	"github.com/JonathonGore/api-check/printer"
	"github.com/JonathonGore/api-check/runner"
)

type RunConfig struct {
	verbose    bool
	standalone bool
}

var (
	rconf = RunConfig{}
)

// Printf prints fmt.Printf statements only run config's verbose setting is true
func printf(text string, args ...interface{}) {
	if rconf.verbose {
		fmt.Printf(text, args...)
	}
}

// Verbose accepts a boolean argument and sets the logging level accordingly
func Verbose(isVerbose bool) {
	rconf.verbose = isVerbose
}

func run(t *testing.T) {
	printf("Running go api-check\n\n")

	// This loads the files containing step definitions below the working directory.
	// Allows the user to have multiple files containing api testing definitions
	files, err := loader.FindRunDefs()
	if err != nil {
		fmt.Printf("Unable to find test definition files: %v\n", err)
		return
	}

	// TODO: It will likely be a good idea to allow the user to specify an alternate
	// config file.
	conf, err := config.New(config.DefaultConfigFile)
	if err != nil {
		fmt.Printf("Unable to parse config file: %v\n", err)
		return
	}

	p := parser.New(conf)

	tests, err := p.Parse(files)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	reports := runner.RunTests(tests)

	if rconf.verbose {
		printer.PrintReports(reports)
	}

	for _, report := range reports {
		if report.Error != nil {
			if t != nil {
				t.Errorf("%v", report.Error) // Signal to go test we have failed a test
			} else {
				// Only exit when running in standalone mode
				os.Exit(1)
			}
		}
	}
}

// RunStandalone configures the runner to use os.Exit(1) if a failure test occurs, and then
// runs the test suite.
func RunStandalone() {
	rconf.standalone = true

	var t *testing.T // Pass in nil testing context
	run(t)
}

// Run reads in *.ac.json files below or in the current directory and runs the test suite
// for each test in every file.
func Run(t *testing.T) {
	if t == nil {
		return // Just return to avoid breaking the users `go test ./... command`
	}

	run(t)
}
