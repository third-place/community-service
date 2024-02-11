package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/third-place/community-service/internal/model"
)

type Post struct {
	gorm.Model
	Text          string
	Draft         bool
	UserID        uint
	User          *User
	Uuid          *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Reports       []*Report  `gorm:"polymorphic:Reported;"`
	Images        []*Image
	Likes         uint
	Replies       uint
	ThreadPostID  uint `gorm:"foreignkey:Post"`
	ThreadPost    *Post
	ReplyToPostID uint `gorm:"foreignkey:Post"`
	ReplyToPost   *Post
	SharePostID   uint `gorm:"foreignkey:Post"`
	SharePost     *Post
}

func (p *Post) GetOwnerUUID() string {
	return p.User.Uuid.String()
}

func CreatePost(user *User, post *model.NewPost) *Post {
	postUuid := uuid.New()
	return &Post{
		Uuid:   &postUuid,
		UserID: user.ID,
		Text:   post.Text,
		Draft:  post.Draft,
	}
}

func CreateShare(user *User, post *Post, share *model.NewShare) *Post {
	postUuid := uuid.New()
	return &Post{
		Uuid:        &postUuid,
		UserID:      user.ID,
		User:        user,
		Text:        share.Text,
		SharePost:   post,
		SharePostID: post.ID,
	}
}

func CreateReply(user *User, post *Post, reply *model.NewReply) *Post {
	return &Post{
		Text:          reply.Text,
		ReplyToPost:   post,
		ReplyToPostID: post.ID,
		UserID:        user.ID,
		User:          user,
	}
}
