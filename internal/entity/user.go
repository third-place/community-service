package entity

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/third-place/community-service/internal/model"
)

type User struct {
	gorm.Model
	Uuid       *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username   string     `gorm:"unique;not null"`
	ProfilePic string
	Name       string
	Role       string           `gorm:"default:'user'"`
	Visibility model.Visibility `gorm:"default:'public'"`
	IsBanned   bool             `gorm:"default:false"`
	Follows    []*Follow
	Posts      []*Post
}

func (u *User) UpdateUserProfileFromModel(user *model.User) {
	u.Name = user.Name
	u.ProfilePic = user.ProfilePic
	u.Username = user.Username
	u.Role = string(user.Role)
	u.IsBanned = user.IsBanned
}
