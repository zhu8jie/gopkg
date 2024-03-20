package xqueue

import (
	"context"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
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
		session.MarkMessage(message, "")
	}
	return nil
}

type KafkaQueue struct {
	addrs    []string
	topics   []string
	groupId  string
	assignor string // sticky, roundrobin, range
	log      *zap.SugaredLogger
}

func NewKafkaQueue(addrs, topics []string, groupId string, assignor string, log *zap.SugaredLogger) *KafkaQueue {
	return &KafkaQueue{
		addrs:    addrs,
		topics:   topics,
		groupId:  groupId,
		assignor: assignor,
		log:      log,
	}
}

func (cfg *KafkaQueue) Start(f func(message *sarama.ConsumerMessage)) (err error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_3_0_0 // 指定 Kafka 版本

	switch cfg.assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	}

	// 创建消费者组
	cg, err := sarama.NewConsumerGroup(cfg.addrs, cfg.groupId, config)
	if err != nil {
		return err
	}

	handler := ConsumerGroupHandler{
		f: f,
	}

	// 开始消费
	go func() {
		for {
			if err := cg.Consume(context.Background(), cfg.topics, handler); err != nil {
				cfg.log.Errorf("consuer topics: %v groupId: %v error %v", cfg.topics, cfg.groupId, err)
			}
		}
	}()

	return nil
}
