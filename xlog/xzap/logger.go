package xzap

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(opts ...Option) *zap.SugaredLogger {
	options := loadOptions(opts...)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		options.Encoder(encoderConfig),
		zapcore.AddSync(options.Out),
		options.Level,
	)

	logger := zap.New(
		core,
		options.ZapOptions...,
	)

	return logger.Sugar()
}

func loadOptions(opts ...Option) *option {
	options := &option{
		Level:      zap.NewAtomicLevelAt(zap.InfoLevel),
		Out:        os.Stdout,
		Encoder:    zapcore.NewConsoleEncoder,
		ZapOptions: []zap.Option{zap.AddCaller(), zap.AddCallerSkip(1)},
	}

	for _, opt := range opts {
		opt.apply(options)
	}

	return options
}
