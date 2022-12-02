package kafka

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/db"
	"github.com/third-place/community-service/internal/mapper"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/repository"
	"log"
)

func InitializeAndRunLoop() {
	userRepository := repository.CreateUserRepository(db.CreateDefaultConnection())
	loopKafkaReader(userRepository)
}

func loopKafkaReader(userRepository *repository.UserRepository) {
	reader, err := GetReader()
	if err != nil {
		return
	}
	log.Print("listening for kafka messages")
	for {
		ev := reader.Poll(60000)
		switch e := ev.(type) {
		case *kafka.Message:
			log.Printf("message received on topic :: %s, data :: %s", e.TopicPartition.String(), string(e.Value))
			if *e.TopicPartition.Topic == "users" {
				readUser(userRepository, e.Value)
			} else if *e.TopicPartition.Topic == "images" {
				updateUserImage(userRepository, e.Value)
			}
		case kafka.Error:
			log.Print("Error :: ", e)
			return
		}
	}
}

func updateUserImage(userRepository *repository.UserRepository, data []byte) {
	result := decodeToMap(data)
	user := result["user"].(map[string]interface{})
	userUuid := user["uuid"].(string)
	s3Key := result["s3_key"].(string)
	log.Print("update user profile pic :: {}, {}, {}", userUuid, s3Key, result)
	userEntity, err := userRepository.FindOneByUuid(uuid.MustParse(userUuid))
	if err != nil {
		log.Print("user not found when updating profile pic")
		return
	}
	log.Print("update user with s3 key", userEntity.Uuid.String(), s3Key)
	userEntity.ProfilePic = s3Key
	userRepository.Save(userEntity)
}

func readUser(userRepository *repository.UserRepository, data []byte) {
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
	userEntity, err := userRepository.FindOneByUuid(uuid.MustParse(userModel.Uuid))
	if err == nil {
		userEntity.UpdateUserProfileFromModel(userModel)
		userRepository.Save(userEntity)
	} else {
		userEntity = mapper.GetUserEntityFromModel(userModel)
		userRepository.Create(userEntity)
	}
}

func decodeToMap(data []byte) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result
}
