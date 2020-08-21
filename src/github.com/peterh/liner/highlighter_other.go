// +build !windows

package liner

import (
	"io"
	"os"
)

const (
	COLOR_RESET = "\x1b[0m"
)

var colorsMap map[Color]string

func initHighlighter(h *Highlighter) {
	colorsMap = map[Color]string{
		COLOR_BLACK:   "\x1b[30m",
		COLOR_BLUE:    "\x1b[34m",
		COLOR_GREEN:   "\x1b[32m",
		COLOR_CYAN:    "\x1b[36m",
		COLOR_RED:     "\x1b[31m",
		COLOR_MAGENTA: "\x1b[35m",
		COLOR_YELLOW:  "\x1b[33m",
		COLOR_WHITE:   "\x1b[37m",
	}
}

func (h *Highlighter) writeColoredOutput(str string, color Color) {
	text := ColorToString(color) + str + COLOR_RESET
	io.WriteString(os.Stdout, text)
}

// Use for non-windows platforms
func ColorToString(c Color) string {
	if str, ok := colorsMap[c]; ok {
		return str
	}
	return ""
}

func isInwinConsole() bool {
	return false
}
