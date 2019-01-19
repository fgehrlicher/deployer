package cmd

import (
	"fmt"
	"github.com/urfave/cli"
	"runtime/debug"
)

var (
	isDebug = false
)

func Init(cliContext *cli.Context) {
	isDebug = cliContext.Parent().Bool("debug")
}

func errorWrap(err error) error {
	if isDebug {
		fmt.Println(string(debug.Stack()))
	}
	return err
}
