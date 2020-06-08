package cli_util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Confirmation string

func NewConfirmation(confirmationQuestion string) *Confirmation {
	confirmation := Confirmation(confirmationQuestion)
	return &confirmation
}

func (this Confirmation) Render() error {
	for {
		fmt.Printf(LF + string(this) + " ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		input = strings.TrimRight(input, LF)
		input = strings.TrimRight(input, CR)
		if input == "y" || input == "yes" {
			break
		} else if input == "n" || input == "no" {
			this.Cleanup()
			fmt.Println(LF + " aborted.")
			os.Exit(0)
		}
		this.Cleanup()
	}
	return nil
}

func (this Confirmation) Cleanup() {
	lineBreaks := strings.Count(string(this), LF) + 2

	for i := 0; i < lineBreaks; i++ {
		fmt.Print(CursorUp + DeleteCurrentLine)
	}
}
