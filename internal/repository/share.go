package repository

import (
	"errors"
	"github.com/third-place/community-service/internal/entity"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type ShareRepository struct {
	conn *gorm.DB
}

func CreateShareRepository(conn *gorm.DB) *ShareRepository {
	return &ShareRepository{conn}
}

func (s *ShareRepository) Save(entity *entity.Post) {
	s.conn.Save(entity)
}

func (s *ShareRepository) FindOneByUuid(shareUuid uuid.UUID) (*entity.Post, error) {
	sharePost := &entity.Post{}
	s.conn.Preload("SharePost").
		Preload("User").
		Where("uuid = ?", shareUuid).
		Find(sharePost)
	if sharePost.ID == 0 {
		return nil, errors.New("no share found")
	}
	return sharePost, nil
}

func (s *ShareRepository) FindByUser(user *entity.User, limit int) []*entity.Post {
	var posts []*entity.Post
	s.conn.
		Preload("User").
		Preload("Images").
		Preload("SharePost").
		Where("user_id = ? AND deleted_at IS NULL AND reply_to_post_id = 0 AND share_post_id != 0", user.ID).
		Order("id desc").
		Limit(limit).
		Find(&posts)
	return posts
}

func (s *ShareRepository) FindByUserFollows(username string, limit int) []*entity.Post {
	var posts []*entity.Post
	s.conn.
		Preload("User").
		Preload("Images").
		Preload("SharePost").
		Table("posts").
		Joins("join follows on follows.following_id = posts.user_id").
		Joins("join users on follows.user_id = users.id").
		Where("users.username = ? AND posts.deleted_at IS NULL AND posts.reply_to_post_id = 0 AND posts.share_post_id != 0", username).
		Order("id desc").
		Limit(limit).
		Find(&posts)
	return posts
}

func (s *ShareRepository) FindAll(limit int) []*entity.Post {
	var posts []*entity.Post
	s.conn.
		Preload("User").
		Preload("Images").
		Preload("SharePost").
		Table("posts").
		Joins("join users on posts.user_id = users.id").
		Where("posts.deleted_at IS NULL AND users.deleted_at IS NULL AND posts.reply_to_post_id = 0 AND posts.share_post_id != 0").
		Order("id desc").
		Limit(limit).
		Find(&posts)
	return posts
}
