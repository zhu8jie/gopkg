package xqueue

import (
	"context"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type KafkaAssignor string

const (
	KAFKA_ASSIGNOR_RANGE      KafkaAssignor = "range"
	KAFKA_ASSIGNOR_STICKY     KafkaAssignor = "sticky"
	KAFKA_ASSIGNOR_ROUNDROBIN KafkaAssignor = "roundrobin"
)

// 定义消费者组处理函数
type ConsumerGroupHandler struct {
	f func(message *sarama.ConsumerMessage)
}

func (h ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// fmt.Printf("Message topic:%q partition:%d offset:%d value:%s\n", message.Topic, message.Partition, message.Offset, string(message.Value))
		h.f(message)
		// session.MarkMessage(message, "")
	}
	return nil
}

type KafkaConsumerGroup struct {
	topics       []string
	logger       *zap.SugaredLogger
	consumerGrop sarama.ConsumerGroup
}

func NewKafkaConsumerGroup(addrs, topics []string, groupId string, assignor KafkaAssignor, log *zap.SugaredLogger) (*KafkaConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_3_0_0 // 指定 Kafka 版本

	switch assignor {
	case KAFKA_ASSIGNOR_STICKY:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case KAFKA_ASSIGNOR_ROUNDROBIN:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case KAFKA_ASSIGNOR_RANGE:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	}

	// 创建消费者组
	cg, err := sarama.NewConsumerGroup(addrs, groupId, config)
	if err != nil {
		return nil, err
	}

	// Track errors
	go func() {
		for err := range cg.Errors() {
			log.Errorf("track error: %v", err)
		}
	}()

	return &KafkaConsumerGroup{
		topics:       topics,
		logger:       log,
		consumerGrop: cg,
	}, nil
}

func (kcg *KafkaConsumerGroup) Start(f func(message *sarama.ConsumerMessage)) {

	handler := ConsumerGroupHandler{
		f: f,
	}

	// 开始消费
	go func() {
		for {
			if err := kcg.consumerGrop.Consume(context.Background(), kcg.topics, handler); err != nil {
				kcg.logger.Errorf("consuer topics: %v error: %v", kcg.topics, err)
			}
		}
	}()
}

func (kcg *KafkaConsumerGroup) Close() error {
	kcg.logger.Debugf("yanw_test %v is closed", kcg.topics)
	return kcg.consumerGrop.Close()
}
