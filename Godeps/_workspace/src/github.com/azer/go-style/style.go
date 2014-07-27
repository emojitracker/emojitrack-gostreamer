package style

import (
	"strings"
	"fmt"
	"regexp"
)

func Style (options, text string) (result string) {
	styles := regexp.MustCompile("[^\\w]+").Split(options, -1)
	result = text

	for _,style := range styles {
		result = fmt.Sprintf("%s%s", ansiCodes[strings.Title(style)], result)
	}

	return result
}

var	ansiCodes = map[string]string{
		"Reset": "\033[0m",
		"Bold": "\033[1m",
		"Italic": "\033[3m",
		"Blink": "\033[5m",
		"Underline": "\033[4m",
		"UnderlineOff": "\033[24m",
		"Inverse": "\033[7m",
		"InverseOff": "\033[27m",
		"Strikethrough": "\033[9m",
		"StrikethroughOff": "\033[29m",

		"Def": "\033[39m",
		"White": "\033[37m",
		"Black": "\033[30m",
		"Grey": "\x1B[90m",
		"Red": "\033[31m",
		"Green": "\033[32m",
		"Blue": "\033[34m",
		"Yellow": "\033[33m",
		"Magenta": "\033[35m",
		"Cyan": "\033[36m",

		"DefBg": "\033[49m",
		"WhiteBg": "\033[47m",
		"BlackBg": "\033[40m",
		"RedBg": "\033[41m",
		"GreenBg": "\033[42m",
		"BlueBg": "\033[44m",
		"YellowBg": "\033[43m",
		"MagentaBg": "\033[45m",
		"CyanBg": "\033[46m",
}
