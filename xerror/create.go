package xerror

import (
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/rabbit-rm/xgo/xstack"
)

const skip = 1

// New 创建一个新的自定义错误，包含堆栈信息
func New(format string, args ...interface{}) error {
	if len(args) == 0 {
		return gerror.NewSkip(skip, format)
	}
	return gerror.NewSkipf(skip, format, args...)
}

func NewWithCaller(format string, args ...interface{}) error {
	if len(args) == 0 {
		return gerror.NewSkip(skip, addCaller(skip, format))
	}
	return gerror.NewSkipf(skip, addCaller(skip, format), args...)
}

// Wrap 包裹其他错误，用于构造多级错误，包含堆栈信息
func Wrap(err error, format string, args ...interface{}) error {
	if len(args) == 0 {
		return gerror.WrapSkip(skip, err, format)
	}
	return gerror.WrapSkipf(skip, err, format, args...)
}

func WrapWithCaller(err error, format string, args ...interface{}) error {
	if len(args) == 0 {
		return gerror.WrapSkip(skip, err, addCaller(skip, format))
	}
	return gerror.WrapSkipf(skip, err, addCaller(skip, format), args...)
}

func addCaller(skip int, format string) string {
	return fmt.Sprintf("%s -> %s", xstack.Caller(skip+2), format)
}
