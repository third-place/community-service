package repository

import (
	"github.com/danielmunro/otto-community-service/internal/entity"
	"github.com/jinzhu/gorm"
)

type ImageRepository struct {
	conn *gorm.DB
}

func CreateImageRepository(conn *gorm.DB) *ImageRepository {
	return &ImageRepository{conn}
}

func (i *ImageRepository) Create(imageEntity *entity.Image) {
	i.conn.Create(imageEntity)
}
