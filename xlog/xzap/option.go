package xzap

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option interface {
	apply(*option)
}

type optionFunc func(*option)

func (f optionFunc) apply(opt *option) {
	f(opt)
}

type option struct {
	Level        zap.AtomicLevel
	Out          io.Writer
	Encoder    func(zapcore.EncoderConfig) zapcore.Encoder
	ZapOptions []zap.Option
}

func WithLevel(level zapcore.Level) Option {
	return optionFunc(func(opt *option) {
		opt.Level = zap.NewAtomicLevelAt(level)
	})
}

func WithOutput(out io.Writer) Option {
	return optionFunc(func(opt *option) {
		opt.Out = out
	})
}

func WithJSONEncoder() Option {
	return optionFunc(func(opt *option) {
		opt.Encoder = zapcore.NewJSONEncoder
	})
}

func EnableCaller() Option {
	return optionFunc(func(opt *option) {
		opt.ZapOptions = append(opt.ZapOptions,zap.AddCaller(),zap.AddCallerSkip(1))
	})
}

func WithZapOptions(zapOpts ...zap.Option) Option {
	return optionFunc(func(opt *option) {
		opt.ZapOptions = append(opt.ZapOptions, zapOpts...)
	})
}
