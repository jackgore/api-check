package suite

import (
	"fmt"
	"os"
	"os/exec"
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

// runScript will execute the bash script in the given filename if non empty.
func runScript(filename string) error {
	if len(filename) == 0 {
		return nil
	}

	cmd := exec.Command("/bin/bash", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if _, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("script %v did not complete sucessfully: %v", filename, err)
	} else if err != nil {
		return fmt.Errorf("unable to execute script: %v", err)
	}

	return nil
}

func run(t *testing.T) error {
	// For now we default the directory we look in for test definitions to be cwd.
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// This loads the files containing step definitions below the working directory.
	// Allows the user to have multiple files containing api testing definitions
	files, err := loader.FindTestDefinitions(dir)
	if err != nil {
		return fmt.Errorf("unable to find test definition files: %v\n", err)
	}

	// TODO: It will likely be a good idea to allow the user to specify an alternate
	// config file.
	conf, err := config.New(config.DefaultConfigFile)
	if err != nil {
		return fmt.Errorf("unable to parse config file: %v\n", err)
	}

	p := parser.New(conf)

	tests, err := p.Parse(files)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := runScript(conf.SetupScript); err != nil {
		return fmt.Errorf("unable to run setup script: %v", err)
	}

	printf("Running go api-check\n\n")
	reports := runner.RunTests(tests)

	if rconf.verbose {
		printer.PrintReports(reports)
	}

	if err := runScript(conf.CleanupScript); err != nil {
		return fmt.Errorf("unable to run cleanup script: %v", err)
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

	return nil
}

// RunStandalone configures the runner to use os.Exit(1) if a failure test occurs, and then
// runs the test suite.
func RunStandalone() {
	rconf.standalone = true

	var t *testing.T // Pass in nil testing context
	if err := run(t); err != nil {
		fmt.Printf("Error running tests: %v\n", err)
		os.Exit(1)
	}
}

// Run reads in *.ac.json files below or in the current directory and runs the test suite
// for each test in every file.
func Run(t *testing.T) {
	if t == nil {
		return // Just return to avoid breaking the users `go test ./... command`
	}

	run(t)
	if err := run(t); err != nil {
		fmt.Printf("Error running tests: %v\n", err)
		os.Exit(1)
	}
}
