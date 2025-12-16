package xkafka

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type SaramaConsumerGroup struct {
	c      sarama.ConsumerGroup
	topics []string
	log    *zap.SugaredLogger
}

func NewSaramaConsumerGroup(addr, topics []string, groupId string, log *zap.SugaredLogger, saramaCfg *sarama.Config) (*SaramaConsumerGroup, error) {

	if saramaCfg == nil {
		// 配置消费者组
		saramaCfg = sarama.NewConfig()
		saramaCfg.Consumer.Return.Errors = true                          // 返回所有错误
		saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest         // 初始偏移量
		saramaCfg.Consumer.Offsets.AutoCommit.Enable = true              // 开启自动提交
		saramaCfg.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second // 自动提交间隔

		// 其他配置
		saramaCfg.Version = sarama.V2_5_0_0
	}

	// 创建消费者组客户端
	client, err := sarama.NewConsumerGroup(addr, groupId, saramaCfg)
	if err != nil {
		return nil, err
	}

	if log == nil {
		l, _ := zap.NewDevelopment()
		log = l.Sugar()
	}

	return &SaramaConsumerGroup{
		c:      client,
		topics: topics,
		log:    log,
	}, nil
}

func (cg *SaramaConsumerGroup) Close() error {
	return cg.c.Close()
}

type GroupSaramaMsg struct {
	Key, Value []byte
	Partition  int32
	Offset     int64
}

type GroupConsumeMsg func(msg GroupSaramaMsg) error

func (cg *SaramaConsumerGroup) Start(f GroupConsumeMsg) {
	cg.log.Debugf("SaramaConsumerGroup start ...")

	// go func() {
	h := ConsumerGroupHandler{
		f:   f,
		log: cg.log,
	}
	err := cg.c.Consume(context.Background(), cg.topics, h)
	if err != nil {
		cg.log.Errorf("SaramaConsumerGroup consume error: %v", err)
	}
	// }()

	cg.log.Debugf("SaramaConsumerGroup end")
}

type ConsumerGroupHandler struct {
	f   GroupConsumeMsg
	log *zap.SugaredLogger
}

// 实现ConsumerGroupHandler接口
func (h ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// fmt.Printf("Message Value: %s\n", string(msg.Value))
		// session.MarkMessage(msg, "")
		err := h.f(GroupSaramaMsg{
			Key:       msg.Key,
			Value:     msg.Value,
			Partition: msg.Partition,
			Offset:    msg.Offset,
		})

		if err != nil {
			h.log.Errorf("ConsumerGroupHandler do message error: %v", err)
		}
	}
	return nil
}
