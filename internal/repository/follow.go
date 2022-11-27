package repository

import (
	"github.com/third-place/community-service/internal/entity"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type FollowRepository struct {
	conn *gorm.DB
}

func CreateFollowRepository(conn *gorm.DB) *FollowRepository {
	return &FollowRepository{conn}
}

func (f *FollowRepository) Create(entity *entity.Follow) {
	f.conn.Create(entity)
}

func (f *FollowRepository) FindByFollowing(user *entity.User) []*entity.Follow {
	var follows []*entity.Follow
	f.conn.Preload("Following").
		Preload("User").
		Where("following_id = ?", user.ID).Find(&follows)
	return follows
}

func (f *FollowRepository) FindByUser(user *entity.User) []*entity.Follow {
	var follows []*entity.Follow
	f.conn.Preload("Following").
		Preload("User").
		Where("user_id = ?", user.ID).Find(&follows)
	return follows
}

func (f *FollowRepository) FindByUserAndFollowing(userUuid uuid.UUID, followingUuid uuid.UUID) *entity.Follow {
	follow := &entity.Follow{}
	f.conn.Raw("SELECT f.* FROM follows f JOIN users u1 ON f.user_id = u1.id JOIN users u2 ON f.following_id = u2.id WHERE u1.uuid = ? AND u2.uuid = ?", userUuid, followingUuid).
		Scan(follow)
	if follow.ID == 0 {
		return nil
	}
	return follow
}

func (f *FollowRepository) FindOne(followUuid uuid.UUID) *entity.Follow {
	follow := &entity.Follow{}
	f.conn.Preload("Following").Where("uuid = ?", followUuid.String()).Find(&follow)
	return follow
}

func (f *FollowRepository) Update(follow *entity.Follow) {
	f.conn.Save(follow)
}
