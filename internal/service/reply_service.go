package service

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	model2 "github.com/third-place/community-service/internal/auth/model"
	"github.com/third-place/community-service/internal/db"
	"github.com/third-place/community-service/internal/entity"
	kafka2 "github.com/third-place/community-service/internal/kafka"
	"github.com/third-place/community-service/internal/mapper"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/repository"
	"log"
)

type ReplyService struct {
	userRepository  *repository.UserRepository
	postRepository  *repository.PostRepository
	replyRepository *repository.ReplyRepository
	securityService *SecurityService
	kafkaWriter     *kafka.Producer
}

func CreateReplyService() *ReplyService {
	conn := db.CreateDefaultConnection()
	return &ReplyService{
		repository.CreateUserRepository(conn),
		repository.CreatePostRepository(conn),
		repository.CreateReplyRepository(conn),
		&SecurityService{},
		kafka2.CreateWriter(),
	}
}

func (r *ReplyService) CreateReply(session *model2.Session, reply *model.NewReply) (*model.Post, error) {
	user, err := r.userRepository.FindOneByUuid(uuid.MustParse(session.User.Uuid))
	if err != nil {
		return nil, err
	}
	post, err := r.postRepository.FindOneByUuid(uuid.MustParse(reply.Post.Uuid))
	if err != nil {
		return nil, err
	}
	replyEntity := entity.CreateReply(user, post, reply)
	r.replyRepository.Create(replyEntity)
	post.Replies += 1
	r.postRepository.Save(post)
	replyModel := mapper.GetPostModelFromEntity(replyEntity)
	err = r.publishPostToKafka(replyModel)
	if err != nil {
		log.Print("error writing reply to kafka :: ", err)
	}
	return replyModel, nil
}

func (r *ReplyService) GetRepliesForPost(postUuid uuid.UUID) ([]*model.Post, error) {
	post, err := r.postRepository.FindOneByUuid(postUuid)
	if err != nil {
		return nil, err
	}
	replies := r.replyRepository.FindRepliesForPost(post)
	return mapper.GetPostModelsFromEntities(replies), nil
}

func (r *ReplyService) publishPostToKafka(post *model.Post) error {
	topic := "posts"
	data, _ := json.Marshal(post)
	return r.kafkaWriter.Produce(
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
