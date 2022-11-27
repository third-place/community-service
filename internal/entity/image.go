package entity

import (
	"github.com/third-place/community-service/internal/model"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Image struct {
	gorm.Model
	S3Key  string
	UserID uint
	User   *User
	Uuid   *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Post   *Post
	PostID uint
	Likes  uint
}

func CreateImage(user *User, post *Post, image *model.NewImage) *Image {
	imageUuid := uuid.MustParse(image.Uuid)
	return &Image{
		UserID: user.ID,
		S3Key:  image.S3Key,
		Uuid:   &imageUuid,
		PostID: post.ID,
	}
}
