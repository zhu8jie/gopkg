package xqueue

import (
	"context"

	"github.com/Shopify/sarama"
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

type KafkaConfig struct {
	Addrs    []string
	Topics   []string
	GroupId  string
	Funcion  func(message *sarama.ConsumerMessage)
	Assignor string // sticky, roundrobin, range
}

func Start(cfg KafkaConfig) (consumerGroup *sarama.ConsumerGroup, err error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_3_0_0 // 指定 Kafka 版本

	switch cfg.Assignor {
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
	cg, err := sarama.NewConsumerGroup(cfg.Addrs, cfg.GroupId, config)
	if err != nil {
		return nil, err
	}

	handler := ConsumerGroupHandler{
		f: cfg.Funcion,
	}

	// 开始消费
	go func() {
		for {
			if err := cg.Consume(context.Background(), cfg.Topics, handler); err != nil {
				panic(err)
			}
		}
	}()

	return &cg, nil
}
