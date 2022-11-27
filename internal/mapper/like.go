package mapper

import (
	"github.com/third-place/community-service/internal/entity"
	"github.com/third-place/community-service/internal/model"
)

func GetPostLikeModelFromEntity(postLike *entity.PostLike) *model.PostLike {
	return &model.PostLike{
		Post: model.Post{
			Uuid: postLike.Post.Uuid.String(),
		},
		User: model.User{
			Uuid: postLike.User.Uuid.String(),
		},
	}
}
