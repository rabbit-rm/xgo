package xrotate

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rabbit-rm/xgo/xerror"
)

var (
	defaultRotateTime = 24 * time.Hour
	defaultMaxAge     = 30 * defaultRotateTime

	// 时间格式对应的最小轮转时间
	timeFormatMinRotation = map[string]time.Duration{
		"%S": time.Second,          // 秒
		"%M": time.Minute,          // 分
		"%H": time.Hour,            // 时
		"%d": 24 * time.Hour,       // 天
		"%m": 28 * 24 * time.Hour,  // 月
		"%y": 365 * 24 * time.Hour, // 年
		"%Y": 365 * 24 * time.Hour, // 年
	}
)

type Option interface {
	apply(*option)
}

type optionFunc func(*option)

func (f optionFunc) apply(opt *option) {
	f(opt)
}

type option struct {
	// LinkName is the name of the symlink to the current log file
	LinkName string
	// RotationTime is the time between rotation
	RotationTime time.Duration
	// MaxAge is the maximum age of a log file before it is removed
	MaxAge time.Duration
	// RotationCount is the maximum number of files to keep
	RotationCount uint
}

// WithLinkName sets the name of the symlink to the current log file
func WithLinkName(linkName string) Option {
	return optionFunc(func(opt *option) {
		opt.LinkName = linkName
	})
}

// WithRotationTime sets the time between rotation
func WithRotationTime(d time.Duration) Option {
	return optionFunc(func(opt *option) {
		opt.RotationTime = d
	})
}

// WithMaxAge sets the maximum age of a log file before it is removed
func WithMaxAge(d time.Duration) Option {
	return optionFunc(func(opt *option) {
		opt.MaxAge = d
	})
}

// WithRotationCount sets the maximum number of files to keep
func WithRotationCount(n uint) Option {
	return optionFunc(func(opt *option) {
		opt.RotationCount = n
	})
}

func loadOptions(opts ...Option) *option {
	options := &option{
		RotationTime: defaultRotateTime,
		MaxAge:       defaultMaxAge,
	}
	for _, opt := range opts {
		opt.apply(options)
	}
	// rotateCount & maxAge 不能同时设置，默认设置 maxAge
	if options.RotationCount != 0 {
		options.MaxAge = 0
	}
	return options
}

// validateRotationTime 校验轮转时间是否与pattern匹配
func validateRotationTime(pattern string, rotationTime time.Duration) error {
	// 找出 pattern 中的最小时间单位
	var minDuration time.Duration = 365 * 24 * time.Hour // 默认为年
	for format, duration := range timeFormatMinRotation {
		if strings.Contains(pattern, format) && duration < minDuration {
			minDuration = duration
		}
	}

	// 如果轮转时间小于 pattern 中的最小时间单位，返回错误
	if rotationTime < minDuration {
		return xerror.Newf("rotation time %v is too small for pattern '%s', minimum allowed is %v",
			rotationTime, pattern, minDuration)
	}

	return nil
}

// NewRotateLogs creates a new rotate logs logger with the given options
func NewRotateLogs(pattern string, options ...Option) (*rotatelogs.RotateLogs, error) {
	opts := loadOptions(options...)

	// 校验轮转时间
	if err := validateRotationTime(pattern, opts.RotationTime); err != nil {
		return nil, err
	}

	return rotatelogs.New(pattern, []rotatelogs.Option{
		rotatelogs.WithRotationCount(opts.RotationCount),
		rotatelogs.WithMaxAge(opts.MaxAge),
		rotatelogs.WithRotationTime(opts.RotationTime),
		rotatelogs.WithLinkName(opts.LinkName),
	}...)
}

// NewRotateLogsFile creates a new rotate logs file with default settings
func NewRotateLogsFile(filename, pattern string, options ...Option) (*rotatelogs.RotateLogs, error) {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return NewRotateLogs(pattern, options...)
}
