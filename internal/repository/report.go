package repository

import (
	"github.com/third-place/community-service/internal/entity"
	"gorm.io/gorm"
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
