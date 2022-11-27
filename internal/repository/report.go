package repository

import (
	"github.com/danielmunro/otto-community-service/internal/entity"
	"github.com/jinzhu/gorm"
)

type ReportRepository struct {
	conn *gorm.DB
}

func CreateReportRepository(conn *gorm.DB) *ReportRepository {
	return &ReportRepository{conn}
}

func (r *ReportRepository) Create(report *entity.Report) {
	r.conn.Create(report)
}
