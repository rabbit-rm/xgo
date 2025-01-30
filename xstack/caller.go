package xstack

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Caller 获得调用方堆栈，skip表示跳过的堆栈
// skip=0 跳过 Caller
func Caller(skip int) string {
	frame := Capture(skip)
	return prettyCaller(frame.File, frame.Line)
}
func prettyCaller(file string, line int) string {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("%s:%d", file, line)
	}
	file = filepath.ToSlash(file)
	dir = filepath.ToSlash(dir)
	file = strings.TrimPrefix(file, dir+"/")
	return fmt.Sprintf("%s:%d", file, line)
}
