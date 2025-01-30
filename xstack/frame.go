package xstack

import (
	"runtime"

	"github.com/rabbit-rm/xgo/internal/stacktrace"
)

func Capture(skip int) runtime.Frame {
	stack := stacktrace.Capture(skip+1, stacktrace.First)
	defer stack.Free()
	frame, _ := stack.Next()
	return frame
}
