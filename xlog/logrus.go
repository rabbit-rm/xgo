//go:build !zap

package xlog

import (
	"github.com/rabbit-rm/xgo/xlog/xlogrus"
	"github.com/sirupsen/logrus"
)

func init() {
	MustSetLogger(&logrusLogger{l: xlogrus.NewLogger()})
}

type logrusLogger struct {
	l *logrus.Logger
}

func (logger *logrusLogger) Debug(args ...interface{}) {
	logger.l.Debug(args...)
}

func (logger *logrusLogger) Info(args ...interface{}) {
	logger.l.Info(args...)
}

func (logger *logrusLogger) Warn(args ...interface{}) {
	logger.l.Warn(args...)
}

func (logger *logrusLogger) Error(args ...interface{}) {
	logger.l.Error(args...)
}

func (logger *logrusLogger) Fatal(args ...interface{}) {
	logger.l.Fatal(args...)
}

func (logger *logrusLogger) Debugf(format string, args ...interface{}) {
	logger.l.Debugf(format, args...)
}

func (logger *logrusLogger) Infof(format string, args ...interface{}) {
	logger.l.Infof(format, args...)
}

func (logger *logrusLogger) Warnf(format string, args ...interface{}) {
	logger.l.Warnf(format, args...)
}

func (logger *logrusLogger) Errorf(format string, args ...interface{}) {
	logger.l.Errorf(format, args...)
}

func (logger *logrusLogger) Fatalf(format string, args ...interface{}) {
	logger.l.Fatalf(format, args...)
}
