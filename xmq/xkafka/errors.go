package xkafka

import (
	"github.com/gogf/gf/v2/errors/gerror"
)

var (
	ErrProducerClosed = gerror.New("producer is closed")
	ErrBrokersEmpty   = gerror.New("brokers is empty")
	ErrGroupIDEmpty   = gerror.New("group ID is empty")
	ErrTopicsEmpty    = gerror.New("topics is empty")
	ErrHandlerNil     = gerror.New("handler is nil")
)
