package entity

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Follow struct {
	gorm.Model
	Uuid        *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID      uint
	User        *User
	Following   *User
	FollowingID uint `gorm:"foreignkey:User"`
}

func GetFollowEntityFromModel(user *User, followUser *User) *Follow {
	return &Follow{
		UserID:      user.ID,
		FollowingID: followUser.ID,
	}
}
