package kafka

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	"log"
)

func InitializeAndRunLoop() {
	loopKafkaReader()
}

func loopKafkaReader() {
	reader, err := GetReader()
	userService := service.CreateUserService()
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
				readUser(userService, e.Value)
			} else if *e.TopicPartition.Topic == "images" {
				updateUserImage(userService, e.Value)
			}
		case kafka.Error:
			log.Print("Error :: ", e)
			return
		}
	}
}

func updateUserImage(userService *service.UserService, data []byte) {
	result := decodeToMap(data)
	user := result["user"].(map[string]interface{})
	userUuid := user["uuid"].(string)
	s3Key := result["s3_key"].(string)
	log.Print("update user profile pic :: {}, {}, {}", userUuid, s3Key, result)
	userModel, err := userService.GetUser(uuid.MustParse(userUuid))
	if err != nil {
		log.Print("user not found when updating profile pic")
		return
	}
	userModel.ProfilePic = s3Key
	userService.UpsertUser(userModel)
}

func readUser(userService *service.UserService, data []byte) {
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
	userService.UpsertUser(userModel)
}

func decodeToMap(data []byte) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result
}
