package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func GetReader() (*kafka.Consumer, error) {
	cfg := createConnectionConfig()
	_ = cfg.SetKey("group.id", "community-service")
	_ = cfg.SetKey("auto.offset.reset", "earliest")
	c, err := kafka.NewConsumer(cfg)
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
