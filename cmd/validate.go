package cmd

import (
	"github.com/urfave/cli"
	"gitlab.osram.info/osram/deployer/git"
	"gitlab.osram.info/osram/deployer/config"
	"fmt"
	"github.com/fatih/color"
	"gitlab.osram.info/osram/deployer/cli_util"
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
