package xerror

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rabbit-rm/xgo/xlog"
)

func TestNew(t *testing.T) {
	err := New("new error")
	fmt.Printf("error: %v\n", err)
	fmt.Printf("error: %+v\n", err)
	xlog.Infof("error: %v", err)
	xlog.Infof("error: %+v", err)
}

func TestNewWithCaller(t *testing.T) {
	err := NewWithCaller("new error with caller")
	fmt.Printf("error: %v\n", err)
	fmt.Printf("error: %+v\n", err)
	xlog.Infof("error: %v", err)
	xlog.Infof("error: %+v", err)
	xlog.Debug("error;%s\n", err.Error())
}

func TestWrap(t *testing.T) {
	err := errors.New("new error")
	err = Wrap(err, "wrap err")
	fmt.Printf("error: %v\n", err)
	fmt.Printf("error: %+v\n", err)
	xlog.Infof("error: %v", err)
	xlog.Infof("error: %+v", err)
}

func TestWrapWithCaller(t *testing.T) {
	err := errors.New("new error")
	err = WrapWithCaller(err, "wrap err with caller")
	fmt.Printf("error: %v\n", err)
	fmt.Printf("error: %+v\n", err)
	xlog.Infof("error: %v", err)
	xlog.Infof("error: %+v", err)
}
