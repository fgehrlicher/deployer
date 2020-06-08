package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli"

	"github.com/fgehrlicher/deployer/cli_util"
	"github.com/fgehrlicher/deployer/config"
	"github.com/fgehrlicher/deployer/git"
)

func ValidateCommand(c *cli.Context) error {
	Init(c)

	conf := config.LoadConfig()
	err := git.InitPackage()
	if err != nil {
		return errorWrap(err)
	}

	_, err = git.NewGitController()
	if err != nil {
		conf.PrintConfig()
		fmt.Println(cli_util.LF + color.RedString(" Configuration invalid.") + cli_util.LF)
		return errorWrap(err)
	}

	conf.PrintConfig()
	fmt.Println(cli_util.LF + color.GreenString(" Configuration valid.") + cli_util.LF)
	return nil
}
