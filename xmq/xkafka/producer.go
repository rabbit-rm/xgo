package xkafka

import (
	"context"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Producer Kafka生产者结构体
type Producer struct {
	sync    sarama.SyncProducer
	async   sarama.AsyncProducer
	config  *sarama.Config
	brokers []string
	mu      sync.RWMutex
	closed  bool
}

// ProducerOption 定义生产者配置选项接口
type ProducerOption interface {
	apply(*Producer)
}

// funcProducerOption 将函数转换为选项接口
type funcProducerOption struct {
	f func(*Producer)
}

func (fpo *funcProducerOption) apply(p *Producer) {
	fpo.f(p)
}

func newFuncProducerOption(f func(*Producer)) *funcProducerOption {
	return &funcProducerOption{f: f}
}

// WithRetry 设置重试次数
func WithRetry(count int) ProducerOption {
	return newFuncProducerOption(func(p *Producer) {
		p.config.Producer.Retry.Max = count
	})
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ProducerOption {
	return newFuncProducerOption(func(p *Producer) {
		p.config.Producer.Timeout = timeout
	})
}

// WithRequiredAcks 设置应答级别
func WithRequiredAcks(acks sarama.RequiredAcks) ProducerOption {
	return newFuncProducerOption(func(p *Producer) {
		p.config.Producer.RequiredAcks = acks
	})
}

// NewProducer 创建新的生产者实例
func NewProducer(brokers []string, opts ...ProducerOption) (*Producer, error) {
	if len(brokers) == 0 {
		return nil, ErrBrokersEmpty
	}

	config := sarama.NewConfig()
	// 设置生产者默认配置
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = 3
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Timeout = 10 * time.Second

	producer := &Producer{
		config:  config,
		brokers: brokers,
	}

	// 应用选项
	for _, opt := range opts {
		opt.apply(producer)
	}

	// 创建同步生产者
	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, gerror.Wrap(err, "create sync producer")
	}
	producer.sync = syncProducer

	// 创建异步生产者
	asyncProducer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		_ = syncProducer.Close()
		return nil, gerror.Wrap(err, "create async producer")
	}
	producer.async = asyncProducer

	// 启动异步生产者的错误处理
	go producer.handleAsyncResults()

	return producer, nil
}

// SendMessage 同步发送消息
func (p *Producer) SendMessage(ctx context.Context, topic string, key, value []byte) error {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return ErrProducerClosed
	}
	p.mu.RUnlock()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	// 使用context控制超时
	done := make(chan error, 1)
	go func() {
		_, _, err := p.sync.SendMessage(msg)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return gerror.Wrap(ctx.Err(), "context canceled")
	case err := <-done:
		if err != nil {
			return gerror.Wrap(err, "send message")
		}
		return nil
	}
}

// SendMessageAsync 异步发送消息
func (p *Producer) SendMessageAsync(topic string, key, value []byte) error {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return ErrProducerClosed
	}
	p.mu.RUnlock()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	p.async.Input() <- msg
	return nil
}

// handleAsyncResults 处理异步生产者的结果
func (p *Producer) handleAsyncResults() {
	for {
		select {
		case success, ok := <-p.async.Successes():
			if !ok {
				return
			}
			_ = success
		case err, ok := <-p.async.Errors():
			if !ok {
				return
			}
			_ = gerror.Wrap(err, "async send message")
		}
	}
}

// Close 关闭生产者
func (p *Producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}
	p.closed = true

	var err error
	if err = p.async.Close(); err != nil {
		err = gerror.Wrap(err, "close async producer")
	}
	if err = p.sync.Close(); err != nil {
		err = gerror.Wrap(err, "close sync producer")
	}

	return err
}
