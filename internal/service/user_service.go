package service

import (
	"github.com/danielmunro/otto-community-service/internal/db"
	"github.com/danielmunro/otto-community-service/internal/entity"
	"github.com/danielmunro/otto-community-service/internal/mapper"
	"github.com/danielmunro/otto-community-service/internal/model"
	"github.com/danielmunro/otto-community-service/internal/repository"
	"github.com/google/uuid"
	"time"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func CreateDefaultUserService() *UserService {
	return CreateUserService(repository.CreateUserRepository(db.CreateDefaultConnection()))
}

func CreateUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository,
	}
}

func (s *UserService) DeleteUser(userUuid uuid.UUID) error {
	userEntity, err := s.userRepository.FindOneByUuid(userUuid)
	if err != nil {
		return err
	}
	deletedAt := time.Now()
	userEntity.DeletedAt = &deletedAt
	s.userRepository.Save(userEntity)
	return nil
}

func (s *UserService) CreateUser(newUser *model.User) *entity.User {
	user := mapper.GetUserEntityFromModel(newUser)
	s.userRepository.Create(user)
	return user
}

func (s *UserService) GetUser(userUuid uuid.UUID) (*model.User, error) {
	userEntity, err := s.userRepository.FindOneByUuid(userUuid)
	if err != nil {
		return nil, err
	}
	return mapper.GetUserModelFromEntity(userEntity), nil
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	userEntity, err := s.userRepository.FindOneByUsername(username)
	if err != nil {
		return nil, err
	}
	return mapper.GetUserModelFromEntity(userEntity), nil
}

func (s *UserService) GetSuggestedFollowsForUser(userUuid uuid.UUID) []*model.User {
	userEntities := s.userRepository.FindUsersNotFollowing(userUuid)
	return mapper.GetUserModelsFromEntities(userEntities)
}
