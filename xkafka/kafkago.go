package xkafka

// import (
// 	"fmt"

// 	kafkago "github.com/confluentinc/confluent-kafka-go/kafka"
// )

// type KafkaGoConf struct {
// 	Topic            string `json:"topic"`
// 	GroupId          string `json:"group.id"`
// 	BootstrapServers string `json:"bootstrap.servers"`
// 	SecurityProtocol string `json:"security.protocol"`
// 	SslCaLocation    string `json:"ssl.ca.location"`
// 	SaslMechanism    string `json:"sasl.mechanism"`
// 	SaslUsername     string `json:"sasl.username"`
// 	SaslPassword     string `json:"sasl.password"`
// }

// func NewKafkaGoProducer(cfg *KafkaGoConf) (*kafkago.Producer, error) {

// 	fmt.Print("init kafka producer, it may take a few seconds to init the connection\n")
// 	//common arguments
// 	var kafkaconf = &kafkago.ConfigMap{
// 		"api.version.request": "true",
// 		"message.max.bytes":   1000000,
// 		"linger.ms":           10,
// 		"retries":             30,
// 		"retry.backoff.ms":    1000,
// 		"acks":                "1"}
// 	kafkaconf.SetKey("bootstrap.servers", cfg.BootstrapServers)

// 	switch cfg.SecurityProtocol {
// 	case "plaintext":
// 		kafkaconf.SetKey("security.protocol", "plaintext")
// 	case "sasl_ssl":
// 		kafkaconf.SetKey("security.protocol", "sasl_ssl")
// 		kafkaconf.SetKey("ssl.ca.location", "conf/ca-cert.pem")
// 		kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
// 		kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
// 		kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)
// 	case "sasl_plaintext":
// 		kafkaconf.SetKey("sasl.mechanism", "PLAIN")
// 		kafkaconf.SetKey("security.protocol", "sasl_plaintext")
// 		kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
// 		kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
// 		kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)
// 	default:
// 		return nil, kafkago.NewError(kafkago.ErrUnknownProtocol, "unknown protocol", true)
// 	}

// 	producer, err := kafkago.NewProducer(kafkaconf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Print("init kafka producer success\n")
// 	return producer, nil
// }
