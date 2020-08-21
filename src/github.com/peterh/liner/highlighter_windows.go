// +build windows

package liner

import (
	"io"
	"os"
	"syscall"
	"unsafe"
)

func initHighlighter(h *Highlighter) {
	h.textAttr = getConsoleTextAttr()
}

func isInwinConsole() bool {
	return true
}

func getConsoleTextAttr() int16 {
	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(syscall.Stdout), uintptr(unsafe.Pointer(&sbi)))

	return sbi.wAttributes
}

func setConsoleTextAttr(attr uint16) (n int, err error) {
	ret, _, err := procSetTextAttribute.Call(uintptr(syscall.Stdout), uintptr(attr))

	// if success, err.Error() is equals "The operation completed successfully."
	if err != nil && err.Error() == "The operation completed successfully." {
		err = nil // set as nil
	}

	return int(ret), err
}

func (h *Highlighter) writeColoredOutput(str string, color Color) {
	setConsoleTextAttr(uint16(color))
	io.WriteString(os.Stdout, str)
	setConsoleTextAttr(uint16(h.textAttr))
}
