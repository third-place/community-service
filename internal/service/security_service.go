package service

import (
	model2 "github.com/third-place/community-service/internal/auth/model"
	"github.com/third-place/community-service/internal/ownable"
)

type SecurityService struct{}

func CreateSecurityService() *SecurityService {
	return &SecurityService{}
}

func (s *SecurityService) Owns(session *model2.Session, ownable ownable.Ownable) bool {
	return session != nil && session.User.Uuid == ownable.GetOwnerUUID()
}
