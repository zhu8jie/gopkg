package xkafka

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"github.com/segmentio/kafka-go"
// 	"log"
// 	"time"
// )

// /*
// People消息的格式 标准情况下应该写在model层的结构体
// */
// type People struct {
// 	Name string `json:"name"`
// 	Pwd  string `json:"pwd"`
// }

// /*
// 正常来说下面的配置host topic partition 应该写在配置文件里
// */
// const host = "localhost:9091"
// const topic = "my"
// const partition = 0

// /*
// NewKafKaCon kafka的客户端连接的初始化方法
// */
// func NewKafKa() (*kafka.Conn, error) {
// 	return kafka.DialLeader(context.Background(), "tcp", host, topic, partition)
// }

// func main() {
// 	writeByConn()
// 	readByConn()

// }

// // writeByConn 基于Conn发送消息
// func writeByConn() {

// 	// 连接至Kafka集群的Leader节点
// 	conn, err := NewKafKaCon()
// 	if err != nil {
// 		log.Fatal("failed to dial leader:", err)
// 	}

// 	// 设置发送消息的超时时间
// 	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
// 	people1 := People{"Tmo",
// 		"124"}
// 	people2 := People{"Mac",
// 		"124"}
// 	people3 := People{"Joker",
// 		"124"}
// 	// 发送消息
// 	str1, _ := json.Marshal(people1)
// 	str2, _ := json.Marshal(people2)
// 	str3, _ := json.Marshal(people3)
// 	_, err = conn.WriteMessages(
// 		kafka.Message{Value: []byte(str1)},
// 		kafka.Message{Value: []byte(str2)},
// 		kafka.Message{Value: []byte(str3)},
// 	)
// 	if err != nil {
// 		log.Fatal("failed to write messages:", err)
// 	}

// 	// 关闭连接
// 	if err := conn.Close(); err != nil {
// 		log.Fatal("failed to close writer:", err)
// 	}
// }

// // readByConn 连接至kafka后接收消息
// func readByConn() {
// 	// 指定要连接的topic和partition

// 	// 连接至Kafka的leader节点
// 	conn, err := NewKafKaCon()
// 	if err != nil {
// 		log.Fatal("failed to dial leader:", err)
// 	}
// 	// 设置读取超时时间
// 	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
// 	// 读取一批消息，得到的batch是一系列消息的迭代器
// 	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

// 	// 遍历读取消息
// 	b := make([]byte, 10e3) // 10KB max per message
// 	for {
// 		p := People{}
// 		n, err := batch.Read(b)
// 		if err != nil {
// 			break
// 		}
// 		err = json.Unmarshal(b[:n], &p)
// 		if err != nil {
// 			fmt.Println(string(b))
// 			fmt.Println(err, "**************")
// 			continue
// 		}
// 		fmt.Println(p)
// 	}
// 	// 关闭batch
// 	if err := batch.Close(); err != nil {
// 		log.Fatal("failed to close batch:", err)
// 	}
// 	// 关闭连接
// 	if err := conn.Close(); err != nil {
// 		log.Fatal("failed to close connection:", err)
// 	}
// }

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/segmentio/kafka-go"
// )

// /*
// People消息的格式 标准情况下应该写在model层的结构体
// */
// type People struct {
// 	Name string `json:"name"`
// 	Pwd  string `json:"pwd"`
// }

// /*
// 正常来说下面的配置host topic partition 应该写在配置文件里
// */
// const host = "localhost:9091"
// const topic = "my"
// const partition = 0

// /*
// NewKafKaCon kafka的客户端连接的初始化方法
// */
// func NewKafKaClient() (*kafka.Conn, error) {
// 	return kafka.DialLeader(context.Background(), "tcp", host, topic, partition)
// }

// func main() {
// 	writeByConn()
// 	readByConn()
// }

// // writeByConn 基于Conn发送消息
// func writeByConn() {

// 	// 连接至Kafka集群的Leader节点
// 	kafkaCli, err := NewKafKaClient()
// 	if err != nil {
// 		log.Fatal("failed to dial leader:", err)
// 	}

// 	// 设置发送消息的超时时间
// 	kafkaCli.SetWriteDeadline(time.Now().Add(10 * time.Second))
// 	people1 := People{"Tmo",
// 		"124"}
// 	people2 := People{"Mac",
// 		"124"}
// 	people3 := People{"Joker",
// 		"124"}
// 	// 发送消息
// 	str1, _ := json.Marshal(people1)
// 	str2, _ := json.Marshal(people2)
// 	str3, _ := json.Marshal(people3)
// 	_, err = kafkaCli.WriteMessages(
// 		kafka.Message{Value: []byte(str1)},
// 		kafka.Message{Value: []byte(str2)},
// 		kafka.Message{Value: []byte(str3)},
// 	)
// 	if err != nil {
// 		log.Fatal("failed to write messages:", err)
// 	}

// 	// 关闭连接
// 	if err := kafkaCli.Close(); err != nil {
// 		log.Fatal("failed to close writer:", err)
// 	}
// }

// func readByConn() {
// 	r := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers:  []string{host},
// 		GroupID:  "consumer-group-id",
// 		Topic:    topic,
// 		MaxBytes: 10e6, // 10MB
// 	})

// 	for {
// 		m, err := r.ReadMessage(context.Background())
// 		if err != nil {
// 			break
// 		}
// 		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
// 	}
// }
