package service

import (
	"github.com/danielmunro/otto-community-service/internal/db"
	"github.com/danielmunro/otto-community-service/internal/repository"
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
