package xkafka

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Consumer Kafka消费者结构体
type Consumer struct {
	client      sarama.ConsumerGroup
	topics      []string
	groupID     string
	handler     Handler
	ready       chan bool
	closeOnce   sync.Once
	ctx         context.Context
	cancel      context.CancelFunc
	config      *sarama.Config
	retryConfig RetryConfig
}

// Handler 定义消息处理器接口
type Handler interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error
}

// ConsumerOption 定义消费者配置选项接口
type ConsumerOption interface {
	apply(*Consumer)
}

// funcConsumerOption 将函数转换为选项接口
type funcConsumerOption struct {
	f func(*Consumer)
}

func (fco *funcConsumerOption) apply(c *Consumer) {
	fco.f(c)
}

func newFuncConsumerOption(f func(*Consumer)) *funcConsumerOption {
	return &funcConsumerOption{f: f}
}

// WithVersion 设置Kafka版本
func WithVersion(version string) ConsumerOption {
	return newFuncConsumerOption(func(c *Consumer) {
		if v, err := sarama.ParseKafkaVersion(version); err == nil {
			c.config.Version = v
		}
	})
}

// WithInitialOffset 设置初始偏移量
func WithInitialOffset(offset int64) ConsumerOption {
	return newFuncConsumerOption(func(c *Consumer) {
		switch offset {
		case sarama.OffsetNewest:
			c.config.Consumer.Offsets.Initial = sarama.OffsetNewest
		case sarama.OffsetOldest:
			c.config.Consumer.Offsets.Initial = sarama.OffsetOldest
		}
	})
}

// WithSessionTimeout 设置会话超时时间
func WithSessionTimeout(timeout time.Duration) ConsumerOption {
	return newFuncConsumerOption(func(c *Consumer) {
		c.config.Consumer.Group.Session.Timeout = timeout
	})
}

// WithHeartbeatInterval 设置心跳间隔
func WithHeartbeatInterval(interval time.Duration) ConsumerOption {
	return newFuncConsumerOption(func(c *Consumer) {
		c.config.Consumer.Group.Heartbeat.Interval = interval
	})
}

// RetryConfig 定义重试配置
type RetryConfig struct {
	MaxRetries     int           // 最大重试次数
	InitialBackoff time.Duration // 初始重试间隔
	MaxBackoff     time.Duration // 最大重试间隔
	Factor         float64       // 重试间隔增长因子
}

// defaultRetryConfig 默认重试配置
var defaultRetryConfig = RetryConfig{
	MaxRetries:     3,
	InitialBackoff: time.Second,
	MaxBackoff:     30 * time.Second,
	Factor:         2.0,
}

// WithRetryConfig 设置重试配置
func WithRetryConfig(config RetryConfig) ConsumerOption {
	return newFuncConsumerOption(func(c *Consumer) {
		c.retryConfig = config
	})
}

// NewConsumer 创建新的消费者实例
func NewConsumer(brokers []string, groupID string, topics []string, handler Handler, opts ...ConsumerOption) (*Consumer, error) {
	if len(brokers) == 0 {
		return nil, ErrBrokersEmpty
	}
	if groupID == "" {
		return nil, ErrGroupIDEmpty
	}
	if len(topics) == 0 {
		return nil, ErrTopicsEmpty
	}
	if handler == nil {
		return nil, ErrHandlerNil
	}

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Session.Timeout = 20 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 6 * time.Second
	config.Consumer.MaxProcessingTime = 500 * time.Millisecond
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_0_0_0

	ctx, cancel := context.WithCancel(context.Background())
	consumer := &Consumer{
		topics:      topics,
		groupID:     groupID,
		handler:     handler,
		ready:       make(chan bool),
		ctx:         ctx,
		cancel:      cancel,
		config:      config,
		retryConfig: defaultRetryConfig,
	}

	// 应用选项
	for _, opt := range opts {
		opt.apply(consumer)
	}

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		cancel()
		return nil, gerror.Wrap(err, "new consumer group")
	}
	consumer.client = client

	return consumer, nil
}

// retryWithBackoff 实现指数退避重试
func (c *Consumer) retryWithBackoff(operation func() error) error {
	var err error
	currentBackoff := c.retryConfig.InitialBackoff

	for retries := 0; retries <= c.retryConfig.MaxRetries; retries++ {
		err = operation()
		if err == nil {
			return nil
		}

		// 如果已经到达最大重试次数，返回最后一次错误
		if retries == c.retryConfig.MaxRetries {
			return gerror.Wrap(err, "max retries reached")
		}

		// 检查是否需要退出
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		case <-time.After(currentBackoff):
			// 计算下一次重试的等待时间
			currentBackoff = time.Duration(float64(currentBackoff) * c.retryConfig.Factor)
			if currentBackoff > c.retryConfig.MaxBackoff {
				currentBackoff = c.retryConfig.MaxBackoff
			}
		}
	}

	return err
}

// Start 启动消费者（增加重试机制）
func (c *Consumer) Start() error {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				err := c.retryWithBackoff(func() error {
					if err := c.client.Consume(c.ctx, c.topics, c.handler); err != nil {
						if errors.Is(err, sarama.ErrClosedConsumerGroup) {
							return err // 不需要重试的错误直接返回
						}
						return gerror.Wrap(err, "consume message failed")
					}
					return nil
				})

				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}

				if err != nil {
					// 记录最终失败的错误
					_ = gerror.Wrap(err, "consumer retry failed")
				}

				// 重置ready通道
				if c.ctx.Err() == nil {
					c.ready = make(chan bool)
				}
			}
		}
	}()

	<-c.ready
	return nil
}

// Stop 停止消费者
func (c *Consumer) Stop() error {
	var err error
	c.closeOnce.Do(func() {
		c.cancel()
		if c.client != nil {
			err = c.client.Close()
		}
	})
	return err
}

// // Setup 实现 ConsumerGroupHandler 接口
// func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
// 	if err := c.handler.Setup(session); err != nil {
// 		return gerror.Wrap(err, "handler setup")
// 	}
// 	close(c.ready)
// 	return nil
// }
//
// // Cleanup 实现 ConsumerGroupHandler 接口
// func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
// 	return c.handler.Cleanup(session)
// }
//
// // ConsumeClaim 实现 ConsumerGroupHandler 接口
// func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
// 	return c.handler.ConsumeClaim(session, claim)
// }
