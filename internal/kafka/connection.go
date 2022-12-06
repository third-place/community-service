package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

func GetReader() (*kafka.Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"security.protocol": os.Getenv("KAFKA_SECURITY_PROTOCOL"),
		"sasl.mechanisms":   os.Getenv("KAFKA_SASL_MECHANISM"),
		"sasl.username":     os.Getenv("KAFKA_SASL_USERNAME"),
		"sasl.password":     os.Getenv("KAFKA_SASL_PASSWORD"),
		"group.id":          "community-service",
		"auto.offset.reset": "smallest",
	})
	if err != nil {
		return nil, err
	}
	err = c.SubscribeTopics([]string{
		"users",
		"images",
	}, nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}
