package mapper

import (
	"github.com/danielmunro/otto-community-service/internal/entity"
	"github.com/danielmunro/otto-community-service/internal/model"
)

func GetImageModelFromEntity(image *entity.Image) *model.Image {
	return &model.Image{
		Uuid:  image.Uuid.String(),
		S3Key: image.S3Key,
	}
}

func GetImageModelsFromEntities(images []*entity.Image) []model.Image {
	imageModels := make([]model.Image, len(images))
	for i, v := range images {
		imageModels[i] = *GetImageModelFromEntity(v)
	}
	return imageModels
}
