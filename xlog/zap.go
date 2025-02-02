//go:build zap

package xlog

import (
	"github.com/rabbit-rm/xgo/xlog/xzap"
	"go.uber.org/zap"
)

type zapLogger struct {
	l *zap.SugaredLogger
}

func init() {
	MustSetLogger(&zapLogger{l: xzap.NewLogger()})
}

func (logger *zapLogger) Debug(args ...interface{}) {
	logger.l.Debug(args...)
}

func (logger *zapLogger) Info(args ...interface{}) {
	logger.l.Info(args...)
}

func (logger *zapLogger) Warn(args ...interface{}) {
	logger.l.Warn(args...)
}

func (logger *zapLogger) Error(args ...interface{}) {
	logger.l.Error(args...)
}

func (logger *zapLogger) Fatal(args ...interface{}) {
	logger.l.Fatal(args...)
}

func (logger *zapLogger) Debugf(format string, args ...interface{}) {
	logger.l.Debugf(format, args...)
}

func (logger *zapLogger) Infof(format string, args ...interface{}) {
	logger.l.Infof(format, args...)
}

func (logger *zapLogger) Warnf(format string, args ...interface{}) {
	logger.l.Warnf(format, args...)
}

func (logger *zapLogger) Errorf(format string, args ...interface{}) {
	logger.l.Errorf(format, args...)
}

func (logger *zapLogger) Fatalf(format string, args ...interface{}) {
	logger.l.Fatalf(format, args...)
}
