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
	"github.com/third-place/community-service/internal/util"
	"log"
)

type ReplyService struct {
	userRepository  *repository.UserRepository
	postRepository  *repository.PostRepository
	replyRepository *repository.ReplyRepository
	securityService *SecurityService
	kafkaProducer   kafka2.Producer
}

func CreateReplyService() *ReplyService {
	conn := db.CreateDefaultConnection()
	return &ReplyService{
		repository.CreateUserRepository(conn),
		repository.CreatePostRepository(conn),
		repository.CreateReplyRepository(conn),
		&SecurityService{},
		kafka2.CreateProducer(),
	}
}

func CreateTestReplyService() *ReplyService {
	conn := util.SetupTestDatabase()
	producer := util.CreateTestProducer()
	return &ReplyService{
		repository.CreateUserRepository(conn),
		repository.CreatePostRepository(conn),
		repository.CreateReplyRepository(conn),
		&SecurityService{},
		producer,
	}
}

func (r *ReplyService) CreateReply(session *model.Session, reply *model.NewReply) (*model.Post, error) {
	user, err := r.userRepository.FindOneInGoodStandingByUuid(uuid.MustParse(session.User.Uuid))
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
	replyModel := model.CreateReply(post.Uuid, user.Uuid, replyEntity.Uuid, replyEntity.Text)
	err = r.PublishReplyToKafka(replyModel)
	if err != nil {
		log.Print("error writing reply to kafka :: ", err)
	}
	return mapper.GetPostModelFromEntity(replyEntity), nil
}

func (r *ReplyService) GetRepliesForPost(postUuid uuid.UUID) ([]*model.Post, error) {
	post, err := r.postRepository.FindOneByUuid(postUuid)
	if err != nil {
		return nil, err
	}
	replies := r.replyRepository.FindRepliesForPost(post)
	return mapper.GetPostModelsFromEntities(replies), nil
}

func (r *ReplyService) GetReply(replyUuid uuid.UUID) (*model.Reply, error) {
	reply, err := r.replyRepository.FindOneByUuid(replyUuid)
	if err != nil {
		return nil, err
	}
	post, err := r.postRepository.FindOneById(reply.ReplyToPostID)
	if err != nil {
		return nil, err
	}
	user, err := r.userRepository.FindOne(post.UserID)
	if err != nil {
		return nil, err
	}
	return model.CreateReply(post.Uuid, user.Uuid, reply.Uuid, reply.Text), nil
}

func (r *ReplyService) PublishReplyToKafka(reply *model.Reply) error {
	topic := "replies"
	data, _ := json.Marshal(reply)
	return r.kafkaProducer.Produce(
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
