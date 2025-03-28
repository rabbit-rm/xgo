package xkafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/IBM/sarama"
)

var brokers = []string{"192.168.157.131:9092"}

func TestNewConsumer(t *testing.T) {
	consumer, err := NewConsumer(brokers, "group-test", []string{"test-topic"}, &MessageHandler{})
	if err != nil {
		return
	}
	defer consumer.Stop()
	// 4. 启动消费者
	if err := consumer.Start(); err != nil {
		return
	}

	// 5. 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("正在关闭消费者...")
}

func TestNewProducer(t *testing.T) {
	producer, err := NewProducer(brokers)
	if err != nil {
		return
	}
	var index int
	for {
		_ = producer.SendMessage(context.Background(), "test-topic", []byte(fmt.Sprintf("test_key_%d", index)), []byte(fmt.Sprintf("test_value_%d", index)))
		index++
		time.Sleep(100 * time.Millisecond)
	}

}

// MessageHandler 实现 kafka.Handler 接口的消息处理器
type MessageHandler struct {
	// 可以在这里添加所需的依赖，比如数据库连接等
}

// Setup 在消费者会话建立时调用
func (h *MessageHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Printf("消费者会话已建立 - 成员ID: %s", session.MemberID())
	return nil
}

// Cleanup 在消费者会话结束时调用
func (h *MessageHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Printf("消费者会话已结束 - 成员ID: %s", session.MemberID())
	return nil
}

// ConsumeClaim 处理消息的主要逻辑
func (h *MessageHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 注意：必须按顺序处理消息，不要启动goroutine
	for msg := range claim.Messages() {
		log.Printf("收到消息: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))

		// 在这里处理消息
		if err := h.processMessage(msg); err != nil {
			log.Printf("处理消息失败: %v", err)
			// 根据业务需求决定是否继续处理
			continue
		}

		// 标记消息已处理
		session.MarkMessage(msg, "")
	}
	return nil
}

// processMessage 处理单条消息的业务逻辑
func (h *MessageHandler) processMessage(msg *sarama.ConsumerMessage) error {
	// 这里实现具体的业务逻辑
	// 比如解析消息、存储数据等
	fmt.Printf("topic:%s,key:%s,value:%s\n", msg.Topic, msg.Key, msg.Value)
	return nil
}
