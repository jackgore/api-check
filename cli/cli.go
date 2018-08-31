package cli

import (
	"github.com/JonathonGore/api-check/suite"
	"github.com/urfave/cli"
)

const (
	confFile         = ".ac.json"
	defaultVerbosity = true
)

// runAction defines the action that is run by invoking `api-check run`
func runAction(c *cli.Context) error {
	suite.Verbose(defaultVerbosity)
	suite.RunStandalone()
	return nil
}

// buildCLICommands builds the list of available commands for the cli.
func buildCLICommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "run",
			Usage:  "runs your test suites",
			Action: runAction,
		},
	}
}

// ConfigureCLI configures a urfave cli app to enable api-check as an cli tool.
func ConfigureCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "api-check"
	app.Usage = "automatically test your APIs"
	app.Version = "0.0.1"
	app.Commands = buildCLICommands()

	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelpAndExit(c, 0)
		return nil
	}

	return app
}
