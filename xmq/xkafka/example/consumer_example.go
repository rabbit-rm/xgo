package example

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"resource/internal/mq/kafka"

	"github.com/IBM/sarama"
)

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
	return nil
}

// StartConsumer 启动消费者的示例函数
func StartConsumer() error {
	// 1. 配置参数
	brokers := []string{"localhost:9092"} // Kafka broker地址
	topics := []string{"your-topic"}      // 要消费的主题
	groupID := "your-group-id"            // 消费者组ID

	// 2. 创建消息处理器
	handler := &MessageHandler{}

	// 3. 创建消费者，使用配置选项
	consumer, err := kafka.NewConsumer(
		brokers,
		groupID,
		topics,
		handler,
		kafka.WithVersion("2.8.0"), // 设置Kafka版本
		kafka.WithInitialOffset(sarama.OffsetOldest), // 从最早的消息开始消费
		kafka.WithSessionTimeout(30*time.Second),     // 设置会话超时时间
	)
	if err != nil {
		return fmt.Errorf("创建消费者失败: %w", err)
	}
	defer consumer.Stop()

	// 4. 启动消费者
	if err := consumer.Start(); err != nil {
		return fmt.Errorf("启动消费者失败: %w", err)
	}

	// 5. 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("正在关闭消费者...")
	return nil
}

// 在main函数或服务启动函数中使用
func main() {
	if err := StartConsumer(); err != nil {
		log.Fatalf("消费者运行失败: %v", err)
	}
}
