//go:build ignore

package xrotate

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

type Option interface {
	apply(*option)
}

type optionFunc func(*option)

func (f optionFunc) apply(opt *option) {
	f(opt)
}

type option struct {
	// Filename is the file to write logs to
	Filename string
	// MaxSize is the maximum size in megabytes of the log file before it gets rotated
	MaxSize int
	// MaxBackups is the maximum number of old log files to retain
	MaxBackups int
	// MaxAge is the maximum number of days to retain old log files
	MaxAge int
	// Compress determines if the rotated log files should be compressed using gzip
	Compress bool
	// LocalTime determines if the time used for formatting the timestamps in backup files is the computer's local time
	LocalTime bool
}

// WithFilename sets the log filename
func WithFilename(filename string) Option {
	return optionFunc(func(opt *option) {
		opt.Filename = filename
	})
}

// WithMaxSize sets the maximum size in megabytes of the log file before it gets rotated
// default 100MB
func WithMaxSize(maxSize int) Option {
	return optionFunc(func(opt *option) {
		opt.MaxSize = maxSize
	})
}

// WithMaxBackups sets the maximum number of old log files to retain
// default 3 days
func WithMaxBackups(maxBackups int) Option {
	return optionFunc(func(opt *option) {
		opt.MaxBackups = maxBackups
	})
}

// WithMaxAge sets the maximum number of days to retain old log files
func WithMaxAge(maxAge int) Option {
	return optionFunc(func(opt *option) {
		opt.MaxAge = maxAge
	})
}

// WithCompress enables or disables compression of rotated files
func WithCompress(compress bool) Option {
	return optionFunc(func(opt *option) {
		opt.Compress = compress
	})
}

// WithLocalTime enables or disables using local time for backup timestamp
func WithLocalTime(localTime bool) Option {
	return optionFunc(func(opt *option) {
		opt.LocalTime = localTime
	})
}

func loadOptions(opts ...Option) *option {
	options := &option{
		Filename:   "app.log",
		MaxSize:    100,  // 100 MB
		MaxBackups: 3,    // keep 3 old files
		MaxAge:     28,   // 28 days
		Compress:   true, // compress rotated files
		LocalTime:  true, // use local time
	}
	for _, opt := range opts {
		opt.apply(options)
	}
	return options
}

// NewLumberjack creates a new lumberjack logger with given options
func NewLumberjack(opts ...Option) *lumberjack.Logger {
	options := loadOptions(opts...)
	return &lumberjack.Logger{
		Filename:   options.Filename,
		MaxSize:    options.MaxSize,
		MaxBackups: options.MaxBackups,
		MaxAge:     options.MaxAge,
		Compress:   options.Compress,
		LocalTime:  options.LocalTime,
	}
}
