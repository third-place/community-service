package service

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	kafka2 "github.com/third-place/community-service/internal/kafka"
	"github.com/third-place/community-service/internal/model"
	"log"
)

type ConsumerService struct {
	userService *UserService
}

func CreateConsumerService() *ConsumerService {
	return &ConsumerService{
		CreateUserService(),
	}
}

func (c *ConsumerService) GetUser(userUuid uuid.UUID) (*model.User, error) {
	return c.userService.GetUser(userUuid)
}

func (c *ConsumerService) UpsertUser(userModel *model.User) {
	c.userService.UpsertUser(userModel)
}

func (c *ConsumerService) Loop() {
	reader, err := kafka2.GetReader()
	if err != nil {
		return
	}
	log.Print("listening for kafka messages")
	for {
		ev := reader.Poll(-1)
		switch e := ev.(type) {
		case *kafka.Message:
			log.Printf("message received on topic :: %s, data :: %s", e.TopicPartition.String(), string(e.Value))
			if *e.TopicPartition.Topic == "users" {
				c.readUser(e.Value)
			} else if *e.TopicPartition.Topic == "images" {
				c.updateUserImage(e.Value)
			}
		case kafka.Error:
			log.Print("Error :: ", e)
			return
		}
	}
}

func (c *ConsumerService) readUser(data []byte) {
	log.Print("consuming user message ", string(data))
	userModel, err := model.DecodeMessageToUser(data)
	if err != nil {
		log.Print("error decoding message to user error :: ", err)
		return
	}
	_, err = uuid.Parse(userModel.Uuid)
	if err != nil {
		return
	}
	c.UpsertUser(userModel)
}

func (c *ConsumerService) updateUserImage(data []byte) {
	result := decodeToMap(data)
	user := result["user"].(map[string]interface{})
	userUuid := user["uuid"].(string)
	key := result["key"].(string)
	log.Print("update user profile pic :: {}, {}, {}", userUuid, key, result)
	userModel, err := c.GetUser(uuid.MustParse(userUuid))
	if err != nil {
		log.Print("user not found when updating profile pic")
		return
	}
	userModel.ProfilePic = key
	c.UpsertUser(userModel)
}

func decodeToMap(data []byte) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result
}
