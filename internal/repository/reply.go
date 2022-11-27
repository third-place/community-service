package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/entity"
)

type ReplyRepository struct {
	conn *gorm.DB
}

func CreateReplyRepository(conn *gorm.DB) *ReplyRepository {
	return &ReplyRepository{conn}
}

func (r *ReplyRepository) Create(reply *entity.Post) {
	r.conn.Create(reply)
}

func (r *ReplyRepository) FindRepliesForPost(post *entity.Post) []*entity.Post {
	var replies []*entity.Post
	r.conn.Preload("User").
		Where("reply_to_post_id = ?", post.ID).
		Order("id desc").
		Find(&replies)
	return replies
}

func (r *ReplyRepository) FindOneByUuid(uuid uuid.UUID) (*entity.Post, error) {
	reply := &entity.Post{}
	r.conn.Where("uuid = ?", uuid).Find(reply)
	if reply.ID == 0 {
		return nil, errors.New(constants.ErrorMessagePostNotFound)
	}
	return reply, nil
}
