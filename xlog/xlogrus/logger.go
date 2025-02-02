package xlogrus

import (
	"bytes"
	"os"

	"github.com/rabbit-rm/xgo/internal/pool"
	"github.com/sirupsen/logrus"
)

func NewLogger(opts ...Option) *logrus.Logger {
	options := loadOptions(opts...)

	// buffer pool
	bfPool := pool.New[*bytes.Buffer](func() *bytes.Buffer {
		return &bytes.Buffer{}
	})

	return &logrus.Logger{
		Out:          options.Out,
		Formatter:    options.Formatter,
		ReportCaller: options.Caller,
		Level:        options.Level,
		BufferPool:   bfPool,
	}
}

func loadOptions(opts ...Option) *option {
	var options = &option{
		Formatter: defaultTextFormatter(),
		Caller:    true,
		Level:     logrus.InfoLevel,
		Out:       os.Stdout,
	}
	for _, opt := range opts {
		opt.apply(options)
	}
	return options
}
