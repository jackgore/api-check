package main

import (
	"fmt"
	"os"

	"github.com/JonathonGore/api-check/config"
	"github.com/JonathonGore/api-check/loader"
	"github.com/JonathonGore/api-check/parser"
	"github.com/JonathonGore/api-check/printer"
	"github.com/JonathonGore/api-check/runner"
)

const (
	confFile = ".ac.json"
)

func main() {
	fmt.Printf("Running go api-check\n\n")

	// This loads the files containing step definitions below the working directory.
	// This allows the user to have multiple files containing api testing definitions
	files, err := loader.FindRunDefs()
	if err != nil {
		fmt.Printf("Unable to find run definition files: %v", err)
		return
	}

	conf, err := config.New(confFile)
	if err != nil {
		fmt.Printf("Unable to parse config object: %v", err)
		return
	}

	p := parser.New(conf)

	tests, err := p.Parse(files)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	reports := runner.RunTests(tests)

	printer.PrintReports(reports)

	for _, report := range reports {
		if report.Error != nil {
			// Exit with non-zero exit code for scripting
			os.Exit(1)
		}
	}
}
