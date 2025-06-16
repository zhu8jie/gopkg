package xkafka

import (
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
)

type XKafkaConsumer struct {
	Addrs    []string
	Topics   []string
	Log      *zap.SugaredLogger
	Consumer sarama.Consumer
}

func NewXKafkaConsumer(addrs, topics []string, log *zap.SugaredLogger) (*XKafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(addrs, config)
	if err != nil {
		return nil, err
	}
	return &XKafkaConsumer{
		Addrs:    addrs,
		Topics:   topics,
		Log:      log,
		Consumer: consumer,
	}, nil
}

func (k *XKafkaConsumer) Close() error {
	err := k.Consumer.Close()
	if err != nil {
		k.Log.Errorf("Consumer.Close error %v:", err)
	}
	return err
}

func (k *XKafkaConsumer) Errorf(format string, args ...interface{}) {
	if k.Log == nil {
		return
	}
	k.Log.Errorf(format, args...)
}
func (k *XKafkaConsumer) Infof(format string, args ...interface{}) {
	if k.Log == nil {
		return
	}
	k.Log.Infof(format, args...)
}
func (k *XKafkaConsumer) Debugf(format string, args ...interface{}) {
	if k.Log == nil {
		return
	}
	k.Log.Debugf(format, args...)
}

type XkafkaMsg struct {
	Key, Value []byte
	Partition  int32
	Offset     int64
}

type XkafkaConsumeMsg func(msg XkafkaMsg) error

func (k *XKafkaConsumer) Start(f XkafkaConsumeMsg) error {
	for _, topic := range k.Topics {
		if topic == "" {
			continue
		}

		partitionList, err := k.Consumer.Partitions(topic)
		if err != nil {
			k.Log.Errorf("consumer.Partitions error: %v", err)
			return err
		}

		for _, partition := range partitionList {
			k.Log.Infof("topic: %v, partition: %v", topic, partition)
			go func(topic string, partition int32) {
				partitionConsumer, err := k.Consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
				if err != nil {
					k.Log.Errorf("loopConsumer err1: %v, partition: %v", err, partition)
					return
				}
				defer partitionConsumer.AsyncClose()

				for {
					msg := <-partitionConsumer.Messages()
					if msg == nil {
						// k.Log.Errorf("partitionConsumer.Messages is nil %v", msg)
						continue
					}
					err := f(XkafkaMsg{
						Key:       msg.Key,
						Value:     msg.Value,
						Partition: msg.Partition,
						Offset:    msg.Offset,
					})

					if err != nil {
						k.Log.Errorf("process func f Get error: %v", err)
					}
				}
			}(topic, partition)
		}
	}
	return nil
}
