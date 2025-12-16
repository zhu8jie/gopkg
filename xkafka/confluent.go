package xkafka

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"sync"
// 	"syscall"

// 	"github.com/Shopify/sarama"
// 	"github.com/confluentinc/confluent-kafka-go/kafka"
// 	"go.uber.org/zap"
// )

// func NewCfConsumerGroup(addr, topics []string, groupId string, log *zap.SugaredLogger, saramaCfg *sarama.Config) (*SaramaConsumerGroup, error) {

// 	if saramaCfg == nil {
// 		// 配置消费者组
// 		saramaCfg = sarama.NewConfig()
// 		saramaCfg.Consumer.Return.Errors = true                  // 返回所有错误
// 		saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest // 从最早的消息开始消费
// 		// saramaCfg.Group.Return.Notifications = true              // 返回通知信息
// 	}

// 	// 创建消费者组客户端
// 	client, err := sarama.NewConsumerGroup(addr, groupId, saramaCfg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if log == nil {
// 		l, _ := zap.NewDevelopment()
// 		log = l.Sugar()
// 	}

// 	return &SaramaConsumerGroup{
// 		c:      client,
// 		topics: topics,
// 		log:    log,
// 	}, nil
// }

// func main() {
// 	// 配置
// 	bootstrapServers := "localhost:9092"
// 	groupID := "my-consumer-group"
// 	topics := []string{"my-topic"}

// 	// 创建消费者配置
// 	config := &kafka.ConfigMap{
// 		"bootstrap.servers":  bootstrapServers,
// 		"group.id":           groupID,
// 		"auto.offset.reset":  "earliest",
// 		"enable.auto.commit": false, // 手动提交偏移量
// 	}

// 	var wg sync.WaitGroup
// 	consumerCount := 3 // 消费者数量

// 	// 创建信号通道
// 	sigchan := make(chan os.Signal, 1)
// 	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

// 	// 创建多个消费者
// 	for i := 0; i < consumerCount; i++ {
// 		wg.Add(1)
// 		go func(consumerID int) {
// 			defer wg.Done()

// 			// 创建消费者实例
// 			consumer, err := kafka.NewConsumer(config)
// 			if err != nil {
// 				log.Printf("Consumer %d: Failed to create consumer: %v", consumerID, err)
// 				return
// 			}
// 			defer consumer.Close()

// 			// 订阅主题
// 			err = consumer.SubscribeTopics(topics, nil)
// 			if err != nil {
// 				log.Printf("Consumer %d: Failed to subscribe: %v", consumerID, err)
// 				return
// 			}

// 			log.Printf("Consumer %d started", consumerID)

// 			run := true
// 			for run {
// 				select {
// 				case <-sigchan:
// 					run = false
// 				default:
// 					// 读取消息，超时1秒
// 					ev := consumer.Poll(1000)
// 					if ev == nil {
// 						continue
// 					}

// 					switch e := ev.(type) {
// 					case *kafka.Message:
// 						// 处理消息
// 						fmt.Printf("Consumer %d: Topic:%s Partition:%d Offset:%d Key:%s Value:%s\n",
// 							consumerID, *e.TopicPartition.Topic, e.TopicPartition.Partition,
// 							e.TopicPartition.Offset, string(e.Key), string(e.Value))

// 						// 手动提交偏移量
// 						_, err := consumer.CommitMessage(e)
// 						if err != nil {
// 							log.Printf("Consumer %d: Commit failed: %v", consumerID, err)
// 						}

// 					case kafka.Error:
// 						// 处理错误
// 						if e.Code() == kafka.ErrAllBrokersDown {
// 							run = false
// 						}
// 						log.Printf("Consumer %d: Error: %v", consumerID, e)
// 					default:
// 						// 忽略其他事件
// 					}
// 				}
// 			}

// 			log.Printf("Consumer %d stopped", consumerID)
// 		}(i)
// 	}

// 	// 等待所有消费者完成
// 	wg.Wait()
// 	log.Println("All consumers stopped")
// }
