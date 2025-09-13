package logger

import (
	"fmt"
	"strconv"
)

const (
	reset = "\033[0m"

	colorCodeRed          = 31
	colorCodeGreen        = 32
	colorCodeYellow       = 33
	colorCodeBlue         = 34
	colorCodeMagenta      = 35
	colorCodeCyan         = 36
	colorCodeLightGray    = 37
	colorCodeDarkGray     = 90
	colorCodeLightRed     = 91
	colorCodeLightGreen   = 92
	colorCodeLightYellow  = 93
	colorCodeLightBlue    = 94
	colorCodeLightMagenta = 95
	colorCodeLightCyan    = 96
	colorCodeWhite        = 97
)

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}
