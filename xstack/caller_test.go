package xstack

import (
	"testing"

	"github.com/rabbit-rm/xgo/xlog"
)

func TestCaller(t *testing.T) {
	caller := Caller(0)
	xlog.Infof("caller: %s", caller)
}

func TestCapture(t *testing.T) {
	frame := Capture(0)
	xlog.Infof("frame file:%s, frame line:%d, frame function:%s", frame.File, frame.Line, frame.Function)
	xlog.Info("")
}
