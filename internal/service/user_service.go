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

type UserService struct {
	userRepository *repository.UserRepository
}

func CreateUserService() *UserService {
	return &UserService{
		repository.CreateUserRepository(db.CreateDefaultConnection()),
	}
}

func CreateTestUserService() *UserService {
	return &UserService{
		repository.CreateUserRepository(util.SetupTestDatabase()),
	}
}

func (s *UserService) DeleteUser(userUuid uuid.UUID) error {
	userEntity, err := s.userRepository.FindOneByUuid(userUuid)
	if err != nil {
		return err
	}
	s.userRepository.Delete(userEntity)
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

func (s *UserService) UpsertUser(userModel *model.User) {
	userEntity, err := s.userRepository.FindOneByUuid(uuid.MustParse(userModel.Uuid))
	if err == nil {
		userEntity.UpdateUserProfileFromModel(userModel)
		s.userRepository.Save(userEntity)
	} else {
		userEntity = mapper.GetUserEntityFromModel(userModel)
		s.userRepository.Create(userEntity)
	}
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
