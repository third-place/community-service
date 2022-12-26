package service

import (
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/db"
	"github.com/third-place/community-service/internal/entity"
	"github.com/third-place/community-service/internal/mapper"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/repository"
	"github.com/third-place/community-service/internal/util"
)

type ReportService struct {
	userRepository   *repository.UserRepository
	postRepository   *repository.PostRepository
	replyRepository  *repository.ReplyRepository
	reportRepository *repository.ReportRepository
}

func CreateReportService() *ReportService {
	conn := db.CreateDefaultConnection()
	return &ReportService{
		repository.CreateUserRepository(conn),
		repository.CreatePostRepository(conn),
		repository.CreateReplyRepository(conn),
		repository.CreateReportRepository(conn),
	}
}

func CreateTestReportService() *ReportService {
	conn := util.SetupTestDatabase()
	return &ReportService{
		repository.CreateUserRepository(conn),
		repository.CreatePostRepository(conn),
		repository.CreateReplyRepository(conn),
		repository.CreateReportRepository(conn),
	}
}

func (r *ReportService) CreatePostReport(newReport *model.NewPostReport) (*model.PostReport, error) {
	user, err := r.userRepository.FindOneInGoodStandingByUuid(uuid.MustParse(newReport.User.Uuid))
	if err != nil {
		return nil, err
	}

	post, err := r.postRepository.FindOneByUuid(uuid.MustParse(newReport.Post.Uuid))
	if err != nil {
		return nil, err
	}

	report := entity.CreateReportPostEntity(user, post, newReport)
	r.reportRepository.Create(report)

	return mapper.GetPostReportModelFromEntity(user, post, report), nil
}

func (r *ReportService) CreateReplyReport(newReport *model.NewPostReport) (*model.PostReport, error) {
	user, err := r.userRepository.FindOneInGoodStandingByUuid(uuid.MustParse(newReport.User.Uuid))
	if err != nil {
		return nil, err
	}

	reply, err := r.replyRepository.FindOneByUuid(uuid.MustParse(newReport.Post.Uuid))
	if err != nil {
		return nil, err
	}

	report := entity.CreateReportPostEntity(user, reply, newReport)
	r.reportRepository.Create(report)

	return mapper.GetPostReportModelFromEntity(user, reply, report), nil
}
