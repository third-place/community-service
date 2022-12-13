package service

import (
	"encoding/json"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/db"
	"github.com/third-place/community-service/internal/entity"
	kafka2 "github.com/third-place/community-service/internal/kafka"
	"github.com/third-place/community-service/internal/mapper"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/repository"
	"github.com/third-place/community-service/internal/util"
	"log"
	"time"
)

type FollowService struct {
	userRepository   *repository.UserRepository
	followRepository *repository.FollowRepository
	kafkaWriter      kafka2.Producer
}

func CreateFollowService() *FollowService {
	conn := db.CreateDefaultConnection()
	return &FollowService{
		repository.CreateUserRepository(conn),
		repository.CreateFollowRepository(conn),
		kafka2.CreateProducer(),
	}
}

func CreateTestFollowService() *FollowService {
	conn := util.SetupTestDatabase()
	producer := util.CreateTestProducer()
	return &FollowService{
		repository.CreateUserRepository(conn),
		repository.CreateFollowRepository(conn),
		producer,
	}
}

func (f *FollowService) CreateFollow(sessionUserUuid uuid.UUID, follow *model.NewFollow) (*model.Follow, error) {
	user, err := f.userRepository.FindOneByUuid(sessionUserUuid)
	if err != nil {
		return nil, err
	}
	toFollow, err := f.userRepository.FindOneByUuid(uuid.MustParse(follow.Following.Uuid))
	if err != nil {
		return nil, err
	}
	followEntity := entity.GetFollowEntityFromModel(user, toFollow)
	f.followRepository.Create(followEntity)
	followModel := mapper.GetFollowModelFromEntity(followEntity, user, toFollow)
	err = f.publishFollowToKafka(followModel)
	if err != nil {
		log.Print("error publishing follow to kafka :: ", err)
	}
	return followModel, nil
}

func (f *FollowService) GetUserFollowers(username string) ([]*model.Follow, error) {
	user, err := f.userRepository.FindOneByUsername(username)
	if err != nil {
		return nil, err
	}

	follows := f.followRepository.FindByFollowing(user)
	return mapper.GetFollowsModelFromEntities(follows), err
}

func (f *FollowService) GetUserFollowersByUsername(username string) ([]*model.Follow, error) {
	user, err := f.userRepository.FindOneByUsername(username)
	if err != nil {
		return nil, err
	}

	follows := f.followRepository.FindByFollowing(user)
	return mapper.GetFollowsModelFromEntities(follows), err
}

func (f *FollowService) GetUserFollows(username string) ([]*model.Follow, error) {
	user, err := f.userRepository.FindOneByUsername(username)
	if err != nil {
		return nil, err
	}

	follows := f.followRepository.FindByUser(user)
	return mapper.GetFollowsModelFromEntities(follows), err
}

func (f *FollowService) GetUserFollowsByUsername(username string) ([]*model.Follow, error) {
	user, err := f.userRepository.FindOneByUsername(username)
	if err != nil {
		return nil, err
	}

	follows := f.followRepository.FindByUser(user)
	return mapper.GetFollowsModelFromEntities(follows), err
}

func (f *FollowService) DeleteFollow(followUuid uuid.UUID, userUuid uuid.UUID) error {
	follow := f.followRepository.FindOne(followUuid)
	if follow == nil {
		log.Print("follow not found :: ", followUuid)
		return errors.New("follow not found")
	}
	user, _ := f.userRepository.FindOneByUuid(userUuid)
	if follow.UserID != user.ID {
		return errors.New("not allowed")
	}
	deletedAt := time.Now()
	follow.DeletedAt = &deletedAt
	f.followRepository.Update(follow)
	return nil
}

func (f *FollowService) publishFollowToKafka(follow *model.Follow) error {
	topic := "follows"
	data, _ := json.Marshal(follow)
	return f.kafkaWriter.Produce(
		&kafka.Message{
			Value: data,
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
		},
		nil,
	)
}
