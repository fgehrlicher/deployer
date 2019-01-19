package cli_util

import (
	"strings"
)

const CursorUp = "\x1b[1A"
const DeleteCurrentLine = "\x1b[2K"

const BoldStart = "\033[1m"
const ColorEnd = "\033[0m"

const CR  = "\r"
const LF  = "\n"
const CRLF = CR+LF

func SanitizeString(inputString string) string {
	inputString = strings.Replace(inputString, LF, "", -1)
	return inputString
}

func Bold(text string) string {
	return BoldStart + text + ColorEnd
}

