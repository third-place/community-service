package service

import (
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/ownable"
)

type SecurityService struct{}

func CreateSecurityService() *SecurityService {
	return &SecurityService{}
}

func (s *SecurityService) Owns(session *model.Session, ownable ownable.Ownable) bool {
	return session != nil && session.User.Uuid == ownable.GetOwnerUUID()
}
