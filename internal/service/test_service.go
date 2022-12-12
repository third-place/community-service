package service

import (
	"github.com/third-place/community-service/internal/entity"
	"github.com/third-place/community-service/internal/model"
)

type TestService struct {
	userService *UserService
	postService *PostService
}

func CreateTestService() *TestService {
	return &TestService{
		userService: CreateTestUserService(),
		postService: CreateTestPostService(),
	}
}

func (t *TestService) CreateUser(user *model.User) *entity.User {
	return t.userService.CreateUser(user)
}

func (t *TestService) CreatePost(session *model.Session, newPost *model.NewPost) (*model.Post, error) {
	return t.postService.CreatePost(session, newPost)
}
