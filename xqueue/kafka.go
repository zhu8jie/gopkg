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
	f      func(message *sarama.ConsumerMessage) error
	logger *zap.SugaredLogger
}

func (h ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// fmt.Printf("Message topic:%q partition:%d offset:%d value:%s\n", message.Topic, message.Partition, message.Offset, string(message.Value))
		err := h.f(message)
		if err != nil {
			h.logger.Errorf("ConsumeClaim error: %v topic: %v message: %v", err, message.Topic, string(message.Value))
			continue
		}
		// 如果对消息不能重复消费，请把标记代码打开，先处理再标记，做到消息不丢失
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
	// config := sarama.NewConfig()
	// config.Consumer.Return.Errors = true
	// config.Version = sarama.V2_3_0_0 // 指定 Kafka 版本

	// switch assignor {
	// case KAFKA_ASSIGNOR_STICKY:
	// 	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	// case KAFKA_ASSIGNOR_ROUNDROBIN:
	// 	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	// case KAFKA_ASSIGNOR_RANGE:
	// 	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	// default:
	// 	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	// }

	// cli, err := sarama.NewClient(addrs, config)
	// if err != nil {
	// 	return nil, err
	// }
	// consumerGroup, err := sarama.NewConsumerGroupFromClient(groupId, cli)
	// if err != nil {
	// 	return nil, err
	// }

	// 创建消费者组
	cg, err := sarama.NewConsumerGroup(addrs, groupId, nil)
	if err != nil {
		return nil, err
	}

	// Track errors
	// go func() {
	// 	for err := range cg.Errors() {
	// 		log.Errorf("track error: %v", err)
	// 	}
	// }()

	return &KafkaConsumerGroup{
		topics:       topics,
		logger:       log,
		consumerGrop: cg,
	}, nil
}

func (kcg *KafkaConsumerGroup) Start(f func(message *sarama.ConsumerMessage) error) {

	// 开始消费
	go func() {
		handler := ConsumerGroupHandler{
			f:      f,
			logger: kcg.logger,
		}

		if err := kcg.consumerGrop.Consume(context.Background(), kcg.topics, handler); err != nil {
			kcg.logger.Errorf("consuer topics: %v error: %v", kcg.topics, err)
		}
	}()
}

func (kcg *KafkaConsumerGroup) Close() error {
	kcg.logger.Debugf("yanw_test %v is closed", kcg.topics)
	return kcg.consumerGrop.Close()
}
