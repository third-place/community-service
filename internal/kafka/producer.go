package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type Producer interface {
	Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error
}

func CreateProducer() Producer {
	producer, err := kafka.NewProducer(createConnectionConfig())
	if err != nil {
		log.Fatal(err)
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
