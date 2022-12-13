package service

import (
	"github.com/third-place/community-service/internal/repository"
)

type ImageService struct {
	imageRepository *repository.ImageRepository
}

func CreateImageService(imageRepository *repository.ImageRepository) *ImageService {
	return &ImageService{
		imageRepository: imageRepository,
	}
}
