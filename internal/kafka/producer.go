package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
)

type Producer interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
}

func CreateProducer() Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		"security.protocol": os.Getenv("KAFKA_SECURITY_PROTOCOL"),
		"sasl.mechanisms":   os.Getenv("KAFKA_SASL_MECHANISM"),
		"sasl.username":     os.Getenv("KAFKA_SASL_USERNAME"),
		"sasl.password":     os.Getenv("KAFKA_SASL_PASSWORD"),
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to create producer: %s", err))
	}
	return producer
}

func CreateMessage(data []byte, topic string) *kafka.Message {
	return &kafka.Message{
		Value: data,
		TopicPartition: kafka.TopicPartition{Topic: &topic,
			Partition: kafka.PartitionAny},
	}
}
