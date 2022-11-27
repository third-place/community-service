package service

import (
	"github.com/third-place/community-service/internal/db"
	"github.com/third-place/community-service/internal/repository"
)

type ImageService struct {
	imageRepository *repository.ImageRepository
}

func CreateDefaultImageService() *ImageService {
	conn := db.CreateDefaultConnection()
	return CreateImageService(repository.CreateImageRepository(conn))
}

func CreateImageService(imageRepository *repository.ImageRepository) *ImageService {
	return &ImageService{
		imageRepository: imageRepository,
	}
}
