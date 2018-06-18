package suite

import (
	"fmt"
	"os"

	"github.com/JonathonGore/api-check/config"
	"github.com/JonathonGore/api-check/loader"
	"github.com/JonathonGore/api-check/parser"
	"github.com/JonathonGore/api-check/printer"
	"github.com/JonathonGore/api-check/runner"
)

type RunConfig struct {
	verbose bool
}

const (
	confFile = ".ac.json"
)

var (
	rconf = RunConfig{}
)

// Printf prints fmt.Printf statements only run config's verbose setting is true
func printf(text string, args ...interface{}) {
	if rconf.verbose {
		fmt.Printf(text, args...)
	}
}

func Verbose(isVerbose bool) {
	rconf.verbose = isVerbose
}

func Run() {
	printf("Running go api-check\n\n")

	// This loads the files containing step definitions below the working directory.
	// Allows the user to have multiple files containing api testing definitions
	files, err := loader.FindRunDefs()
	if err != nil {
		printf("Unable to find test definition files: %v\n", err)
		return
	}

	conf, err := config.New(confFile)
	if err != nil {
		printf("Unable to parse config file: %v\n", err)
		return
	}

	p := parser.New(conf)

	tests, err := p.Parse(files)
	if err != nil {
		printf("%v\n", err)
		return
	}

	reports := runner.RunTests(tests)

	if rconf.verbose {
		printer.PrintReports(reports)
	}

	for _, report := range reports {
		if report.Error != nil {
			// Exit with non-zero exit code for scripting
			os.Exit(1)
		}
	}
}
