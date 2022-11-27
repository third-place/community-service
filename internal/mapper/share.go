package mapper

import (
	"github.com/danielmunro/otto-community-service/internal/entity"
	"github.com/danielmunro/otto-community-service/internal/model"
)

func GetShareModelFromEntity(post *entity.Post) *model.Share {
	return &model.Share{
		Post:      *GetPostModelFromEntity(post.SharePost),
		Uuid:      post.Uuid.String(),
		Text:      post.Text,
		CreatedAt: post.CreatedAt,
	}
}

func GetShareModelsFromEntities(posts []*entity.Post) []*model.Share {
	shareModels := make([]*model.Share, len(posts))
	for i, v := range posts {
		shareModels[i] = GetShareModelFromEntity(v)
	}
	return shareModels
}
