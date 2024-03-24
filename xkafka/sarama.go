package xkafka

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
	"github.com/zhu8jie/gopkg/xutils"
	"go.uber.org/zap"
)

type SaramaConf struct {
	Addrs []string
	Topic string
	Log   *zap.SugaredLogger
}

type SaramaProducer struct {
	producer sarama.AsyncProducer
	addrs    []string
	topic    string
	log      *zap.SugaredLogger
}

func NewSaramaProducer(conf *SaramaConf) (*SaramaProducer, error) {
	ret := new(SaramaProducer)
	ret.addrs = conf.Addrs
	ret.topic = conf.Topic
	ret.log = conf.Log

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionGZIP     // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	config.Version = sarama.V2_2_0_0

	// config := sarama.NewConfig()
	// config.Producer.RequiredAcks = sarama.WaitForLocal        // 发送完数据需要leader和follow都确认
	// config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	// config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 连接kafka
	p, err := sarama.NewAsyncProducer(conf.Addrs, config)
	ret.producer = p

	// sarama.Logger =

	if conf.Log != nil {
		go func() {
			for err := range p.Errors() {
				conf.Log.Errorf("AsyncProducer get error: %v", err)
			}
		}()
	}

	return ret, err
}

func (p *SaramaProducer) SendMessage(message string) {
	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Value:     sarama.StringEncoder(message),
		Timestamp: time.Now(),
	}

	p.producer.Input() <- msg
}

func (p *SaramaProducer) Close() {
	p.producer.AsyncClose()
}

type SaramaConsumer struct {
	consumer sarama.Consumer
	addrs    []string
	topic    string
	log      *zap.SugaredLogger
}

func NewSaramaConsumer(conf *SaramaConf) (*SaramaConsumer, error) {
	ret := new(SaramaConsumer)
	ret.addrs = conf.Addrs
	ret.topic = conf.Topic
	ret.log = conf.Log

	config := sarama.NewConfig()
	config.Version = sarama.V2_2_0_0

	// 连接kafka
	consumer, err := sarama.NewConsumer(conf.Addrs, config)
	ret.consumer = consumer

	return ret, err
}

type SaramaMsg struct {
	Key, Value []byte
	Partition  int32
	Offset     int64
}

type ConsumeMsg func(msg SaramaMsg) error

func (sc *SaramaConsumer) Start(f ConsumeMsg) error {
	sc.log.Debugf("SaramaConsumer start ...")
	partitionList, err := sc.consumer.Partitions(sc.topic) // 根据topic取到所有的分区
	if err != nil {
		sc.log.Errorf("fail to get list of partition:err%v\n", err)
		return err
	}

	for partition := range partitionList { // 遍历所有的分区
		// 异步从每个分区消费信息
		go func(partition int) {
			// 针对每个分区创建一个对应的分区消费者
			pc, err := sc.consumer.ConsumePartition(sc.topic, int32(partition), sarama.OffsetNewest)
			if err != nil {
				sc.log.Errorf("failed to start consumer for partition %d,err:%v\n", partition, err)
			}
			defer pc.AsyncClose()
			// sc.log.Debugf("SaramaConsumer pc number: %v", partition)

			for msg := range pc.Messages() {
				// sc.log.Debugf("SaramaConsumer get message: %v", msg)
				err := f(SaramaMsg{
					Key:       msg.Key,
					Value:     msg.Value,
					Partition: msg.Partition,
					Offset:    msg.Offset,
				})
				if err != nil {
					sc.log.Errorf("SaramaConsumer consume error: %v", err)
				}
			}
		}(partition)
	}
	return nil
}

func (sc *SaramaConsumer) Close() error {
	err := sc.consumer.Close()
	if err != nil {
		return err
	}

	return nil
}

type ConsumerGroupConf struct {
	Addrs   []string
	Topic   []string
	Log     *zap.SugaredLogger
	GroupId string
}

type SaramaConsumerGroup struct {
	consumer sarama.ConsumerGroup
	Topic    []string
	log      *zap.SugaredLogger
}

func NewSaramaConsumerGroup(conf *ConsumerGroupConf) (*SaramaConsumerGroup, error) {
	ret := new(SaramaConsumerGroup)
	ret.Topic = conf.Topic
	ret.log = conf.Log

	config := sarama.NewConfig()
	config.Version = sarama.V2_2_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky

	if conf.GroupId == "" {
		conf.GroupId = xutils.Int64ToStr(time.Now().UnixNano())
	}
	// 连接kafka
	consumerGroup, err := sarama.NewConsumerGroup(conf.Addrs, conf.GroupId, config)
	ret.consumer = consumerGroup

	return ret, err
}

func (sc *SaramaConsumerGroup) Start(f ConsumeMsg) error {
	sc.log.Debugf("SaramaConsumerGroup start ...")

	go func() {
		h := ConsumerGroupHandler{
			f:   f,
			log: sc.log,
		}
		err := sc.consumer.Consume(context.Background(), sc.Topic, h)
		if err != nil {
			sc.log.Errorf("SaramaConsumerGroup consume error: %v", err)
		}
	}()

	// sc.log.Debugf("SaramaConsumer partition for end")
	return nil
}

func (sc *SaramaConsumerGroup) Close() error {
	err := sc.consumer.Close()
	if err != nil {
		return err
	}

	return nil
}

type ConsumerGroupHandler struct {
	f   ConsumeMsg
	log *zap.SugaredLogger
}

// 实现ConsumerGroupHandler接口
func (h ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// fmt.Printf("Message Value: %s\n", string(msg.Value))
		// session.MarkMessage(msg, "")
		err := h.f(SaramaMsg{
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
