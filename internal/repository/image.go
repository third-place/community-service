package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/third-place/community-service/internal/entity"
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
