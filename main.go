package main

import (
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"github.com/fgehrlicher/deployer/cmd"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Usage = app.Name + " command [command options]"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug mode",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "validate",
			Aliases: []string{"v"},
			Usage:   "validates the current settings",
			Action:  cmd.ValidateCommand,
		},
		{
			Name:    "manual",
			Aliases: []string{"m"},
			Usage:   "shows the deployment manual",
			Action:  cmd.ManualCommand,
		},
		{
			Name:    "deploy",
			Aliases: []string{"d"},
			Usage:   "marks an commit for deployment",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "latest, l",
					Usage: "select the latest availabe commit/tag",
				},
				cli.StringFlag{
					Name:  "stage, s",
					Usage: "the target `STAGE` (dev, qa, prod)",
				},
				cli.StringFlag{
					Name:  "type, t",
					Usage: "the version `TYPE` (increment, patch, minor, major)",
				},
			},
			Action: cmd.DeployCommand,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
}
