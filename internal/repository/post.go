package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/entity"
	"gorm.io/gorm"
)

type PostRepository struct {
	conn *gorm.DB
}

func CreatePostRepository(conn *gorm.DB) *PostRepository {
	return &PostRepository{conn}
}

func (p *PostRepository) Create(entity *entity.Post) {
	p.conn.Create(entity)
}

func (p *PostRepository) Save(entity *entity.Post) {
	p.conn.Save(entity)
}

func (p *PostRepository) Delete(entity *entity.Post) {
	p.conn.Delete(entity)
}

func (p *PostRepository) FindDraftsByUser(user *entity.User, limit int) []*entity.Post {
	return p.FindByUser(user, limit, true)
}

func (p *PostRepository) FindPublishedByUser(user *entity.User, limit int) []*entity.Post {
	return p.FindByUser(user, limit, false)
}

func (p *PostRepository) FindByUser(user *entity.User, limit int, draft bool) []*entity.Post {
	var posts []*entity.Post
	p.conn.
		Preload("User").
		Preload("Images").
		Preload("SharePost").
		Table("posts").
		Joins("JOIN users ON posts.user_id = users.id").
		Where("posts.user_id = ? AND users.is_banned = false AND posts.deleted_at IS NULL AND posts.reply_to_post_id IS NULL AND posts.draft = ?", user.ID, draft).
		Order("posts.id desc").
		Limit(limit).
		Find(&posts)
	return posts
}

func (p *PostRepository) FindByLikes(user *entity.User, limit int) []*entity.Post {
	var posts []*entity.Post
	p.conn.
		Preload("User").
		Preload("Images").
		Preload("SharePost").
		Table("posts").
		Joins("join post_likes on post_likes.post_id = posts.id").
		Joins("JOIN users ON posts.user_id = users.id").
		Where("post_likes.user_id = ? AND users.is_banned = false AND posts.deleted_at IS NULL AND posts.reply_to_post_id IS NULL AND posts.draft = false", user.ID).
		Order("posts.id desc").
		Limit(limit).
		Find(&posts)
	return posts
}

func (p *PostRepository) FindOneByUuid(uuid uuid.UUID) (*entity.Post, error) {
	post := &entity.Post{}
	p.conn.
		Preload("User").
		Preload("Images").
		Preload("SharePost").
		Table("posts").
		Joins("JOIN users ON posts.user_id = users.id").
		Where("posts.uuid = ? AND users.is_banned = false AND posts.deleted_at IS NULL", uuid).
		Find(post)
	if post.ID == 0 {
		return nil, errors.New(constants.ErrorMessagePostNotFound)
	}
	return post, nil
}

func (p *PostRepository) FindOneById(id uint) (*entity.Post, error) {
	post := &entity.Post{}
	p.conn.
		Preload("User").
		Preload("Images").
		Preload("SharePost").
		Table("posts").
		Joins("JOIN users ON posts.user_id = users.id").
		Where("posts.id = ? AND users.is_banned = false AND posts.deleted_at IS NULL", id).
		Find(post)
	if post.ID == 0 {
		return nil, errors.New(constants.ErrorMessagePostNotFound)
	}
	return post, nil
}

func (p *PostRepository) FindByUserFollows(username string, limit int) []*entity.Post {
	var posts []*entity.Post
	p.conn.
		Preload("User").
		Preload("Images").
		Preload("SharePost").
		Table("posts").
		Joins("join follows on follows.following_id = posts.user_id").
		Joins("join users on follows.user_id = users.id").
		Where("users.username = ? AND users.is_banned = false AND posts.deleted_at IS NULL AND posts.reply_to_post_id IS NULL AND posts.draft = false", username).
		Order("id desc").
		Limit(limit).
		Find(&posts)
	return posts
}

func (p *PostRepository) FindAll(limit int) []*entity.Post {
	var posts []*entity.Post
	p.conn.
		Preload("User").
		Preload("Images").
		Preload("SharePost").
		Table("posts").
		Joins("join users on posts.user_id = users.id").
		Where("posts.deleted_at IS NULL AND users.is_banned = false AND users.deleted_at IS NULL AND posts.reply_to_post_id = 0 AND posts.draft = false").
		Order("id desc").
		Limit(limit).
		Find(&posts)
	return posts
}

func (p *PostRepository) FindByIDs(ids []uint) []*entity.Post {
	var posts []*entity.Post
	p.conn.
		Preload("User").
		Preload("Images").
		Preload("SharePost").
		Table("posts").
		Joins("join users on posts.user_id = users.id").
		Where("posts.deleted_at IS NULL AND users.is_banned = false AND posts.id IN (?) AND posts.draft = false", ids).
		Order("id desc").
		Find(&posts)
	return posts
}
