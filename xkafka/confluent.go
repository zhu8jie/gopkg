package xkafka

import (
	"math/rand"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/zhu8jie/gopkg/xutils"
	"go.uber.org/zap"
)

type ConfluentConsumer struct {
	Consumer *kafka.Consumer
	Logger   *zap.SugaredLogger
}

func NewKafkaConsumer(addrs, topics []string, group string, log *zap.SugaredLogger) (*ConfluentConsumer, error) {
	if group == "" {
		group = "my-group-" + xutils.IntToStr(rand.Intn(100))
	}
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(addrs, ","), // Kafka服务器地址
		"group.id":          group,                    // 消费者组ID
		"auto.offset.reset": "earliest",               // 当没有偏移量时从哪里开始消费
	})
	if err != nil {
		return nil, err
	}

	// 订阅主题
	c.SubscribeTopics(topics, nil)

	return &ConfluentConsumer{
		Consumer: c,
		Logger:   log,
	}, nil
}

func (lc *ConfluentConsumer) Close() {
	lc.Consumer.Close()
}

type ConfluentConsumerMsg func(msgStr string) error

func (lc *ConfluentConsumer) Start(f ConfluentConsumerMsg) {
	// 持续消费消息
	for {
		msg, err := lc.Consumer.ReadMessage(-1) // -1 表示无限等待，直到有消息到来或发生错误
		if err != nil {
			lc.Logger.Errorf("LibKafka_consumer_ReadMessage_error: %v", err)
			break
		}

		// fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		err = f(msg.String())
		if err != nil {
			lc.Logger.Errorf("LibKafka_consumer_do_error: %v", err)
		}

	}
	lc.Consumer.Close()
}
