package service

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/db"
	"github.com/third-place/community-service/internal/entity"
	kafka2 "github.com/third-place/community-service/internal/kafka"
	"github.com/third-place/community-service/internal/mapper"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/repository"
)

type LikeService struct {
	likeRepository *repository.LikeRepository
	postRepository *repository.PostRepository
	userRepository *repository.UserRepository
	kafkaWriter    kafka2.Producer
}

func CreateDefaultLikeService() *LikeService {
	conn := db.CreateDefaultConnection()
	return &LikeService{
		likeRepository: repository.CreateLikeRepository(conn),
		postRepository: repository.CreatePostRepository(conn),
		userRepository: repository.CreateUserRepository(conn),
		kafkaWriter:    kafka2.CreateProducer(),
	}
}

func (l *LikeService) CreateLikeForPost(postUuid uuid.UUID, userUuid uuid.UUID) (*model.PostLike, error) {
	user, err := l.userRepository.FindOneInGoodStandingByUuid(userUuid)
	if err != nil {
		return nil, err
	}
	post, err := l.postRepository.FindOneByUuid(postUuid)
	if err != nil {
		return nil, err
	}
	newPostLike := &entity.PostLike{
		Post: post,
		User: user,
	}
	l.likeRepository.Create(newPostLike)
	postModel := mapper.GetPostLikeModelFromEntity(newPostLike)
	err = l.publishPostLikeToKafka(postModel)
	return postModel, err
}

func (l *LikeService) DeleteLikeForPost(postUuid uuid.UUID, userUuid uuid.UUID) error {
	user, err := l.userRepository.FindOneInGoodStandingByUuid(userUuid)
	if err != nil {
		return err
	}
	post, err := l.postRepository.FindOneByUuid(postUuid)
	if err != nil {
		return err
	}
	postLike, err := l.likeRepository.FindByPostAndUser(post, user)
	if err != nil {
		return nil
	}
	l.likeRepository.DeletePostLike(postLike)
	return nil
}

func (l *LikeService) publishPostLikeToKafka(like *model.PostLike) error {
	topic := "post-likes"
	data, _ := json.Marshal(like)
	return l.kafkaWriter.Produce(
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
