package xlogrus

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rabbit-rm/xgo/internal/pkg"
	"github.com/rabbit-rm/xgo/internal/stacktrace"
	"github.com/sirupsen/logrus"
)

type Option interface {
	apply(*option)
}

type optionFunc func(*option)

func (f optionFunc) apply(opt *option) {
	f(opt)
}

// Logrus Options
type option struct {
	Formatter logrus.Formatter
	Caller    bool
	Level     logrus.Level
	Out          io.Writer
}

func WithFormatter(formatter logrus.Formatter) Option {
	return optionFunc(func(opt *option) {
		opt.Formatter = formatter
	})
}

func EnableReportCaller() Option {
	return optionFunc(func(opt *option) {
		opt.Caller = true
	})
}

func WithLevel(level logrus.Level) Option {
	return optionFunc(func(opt *option) {
		opt.Level = level
	})
}

func WithOut(out io.Writer) Option {
	return optionFunc(func(opt *option) {
		opt.Out = out
	})
}

const defaultTimeFormatLayer = "2006-01-02T15:04:05"

func defaultTextFormatter() *logrus.TextFormatter {
	return &logrus.TextFormatter{
		ForceQuote:       true,
		TimestampFormat:  defaultTimeFormatLayer,
		CallerPrettyfier: callerPretty,
	}
}

func defaultJsonFormatter() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat:  defaultTimeFormatLayer,
		CallerPrettyfier: callerPretty,
	}
}

// 自定义调用堆栈输出
func callerPretty(_ *runtime.Frame) (function string, file string) {
	frame := getCaller()
	file, line := frame.File, frame.Line
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Sprintf("%s:%d", file, line)
	}
	file = filepath.ToSlash(file)
	dir = filepath.ToSlash(dir)
	file = strings.TrimPrefix(file, dir+"/")
	return "", fmt.Sprintf("%s:%d", file, line)
}

var (
	logrusDepth = 8
)

func getCaller() *runtime.Frame {
	var f *runtime.Frame
	stacks := stacktrace.Capture(logrusDepth+2, stacktrace.Full)
	for frame, more := stacks.Next(); more; frame, more = stacks.Next() {
		pkgName := getPkgName(frame.Function)
		f = &frame
		if strings.HasPrefix(pkgName, pkg.LogrusName()) ||
			strings.HasPrefix(pkgName, pkg.XLogName()) {
		} else {
			break
		}
	}
	return f
}

func getPkgName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}
