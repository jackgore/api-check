package main

import (
	"log"
	"os"

	"github.com/JonathonGore/api-check/suite"
	"github.com/urfave/cli"
)

const (
	confFile         = ".ac.json"
	defaultVerbosity = true
)

func ConfigureCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "api-check"
	app.Usage = "automatically test your APIs"
	app.Version = "0.0.1"

	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			cli.ShowAppHelpAndExit(c, 0)
		}

		command := c.Args().Get(0)
		if command == "run" {
			suite.Verbose(defaultVerbosity)
			suite.RunStandalone()
		} else {
			cli.ShowAppHelpAndExit(c, 0)
		}

		return nil
	}

	return app
}

func main() {
	app := ConfigureCLI()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
