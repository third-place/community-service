package service

import (
	model2 "github.com/danielmunro/otto-community-service/internal/auth/model"
	"github.com/danielmunro/otto-community-service/internal/db"
	"github.com/danielmunro/otto-community-service/internal/entity"
	"github.com/danielmunro/otto-community-service/internal/mapper"
	"github.com/danielmunro/otto-community-service/internal/model"
	"github.com/danielmunro/otto-community-service/internal/repository"
	"github.com/google/uuid"
)

type ReplyService struct {
	userRepository  *repository.UserRepository
	postRepository  *repository.PostRepository
	replyRepository *repository.ReplyRepository
	securityService *SecurityService
}

func CreateReplyService() *ReplyService {
	conn := db.CreateDefaultConnection()
	return &ReplyService{
		repository.CreateUserRepository(conn),
		repository.CreatePostRepository(conn),
		repository.CreateReplyRepository(conn),
		&SecurityService{},
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
