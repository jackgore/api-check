package cli

import (
	"fmt"

	"github.com/JonathonGore/api-check/config"
	"github.com/JonathonGore/api-check/parser"
	"github.com/JonathonGore/api-check/suite"
	"github.com/urfave/cli"
)

const (
	confFile         = ".ac.json"
	defaultVerbosity = true
)

// verifyAction defines the action that is run by incoking `api-check verify <filename>`
func verifyAction(c *cli.Context) error {
	if c.NArg() == 0 {
		return cli.NewExitError(fmt.Sprintf("%v %v requires at least 1 argument", c.App.Name, c.Command.Name), 1)
	}

	conf, err := config.New(config.DefaultConfigFile)
	if err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
	}

	p := parser.New(conf)

	if _, err := p.Parse([]string(c.Args())); err != nil {
		return cli.NewExitError(fmt.Sprintf("%v", err), 1)
	}

	fmt.Println("validated successfuly")
	return nil
}

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
		{
			Name:   "verify",
			Usage:  "verify an api-check file",
			Action: verifyAction,
		},
	}
}

// ConfigureCLI configures a urfave cli app to enable api-check as an cli tool.
func ConfigureCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "api-check"
	app.Usage = "automatically test your APIs"
	app.Version = "0.0.1"
	//	app.ExitErrHandler = handleError
	app.Commands = buildCLICommands()

	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelpAndExit(c, 0)
		return nil
	}

	return app
}