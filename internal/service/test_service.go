package service

import (
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/entity"
	"github.com/third-place/community-service/internal/model"
)

type TestService struct {
	userService   *UserService
	postService   *PostService
	followService *FollowService
}

func CreateTestService() *TestService {
	return &TestService{
		CreateTestUserService(),
		CreateTestPostService(),
		CreateTestFollowService(),
	}
}

func (t *TestService) CreateUser(user *model.User) *entity.User {
	return t.userService.CreateUser(user)
}

func (t *TestService) CreatePost(session *model.Session, newPost *model.NewPost) (*model.Post, error) {
	return t.postService.CreatePost(session, newPost)
}

func (t *TestService) GetPost(session *model.Session, postUuid uuid.UUID) (*model.Post, error) {
	return t.postService.GetPost(session, postUuid)
}

func (t *TestService) DeletePost(session *model.Session, postUuid uuid.UUID) error {
	return t.postService.DeletePost(session, postUuid)
}

func (t *TestService) GetPostsFirehose(username *string, limit int) ([]*model.Post, error) {
	return t.postService.GetPostsFirehose(username, limit)
}

func (t *TestService) GetPostsForUser(username string, viewerUuid *uuid.UUID, limit int) ([]*model.Post, error) {
	return t.postService.GetPostsForUser(username, viewerUuid, limit)
}

func (t *TestService) CreateFollow(sessionUserUuid uuid.UUID, followUuid uuid.UUID) (*model.Follow, error) {
	return t.followService.CreateFollow(
		sessionUserUuid,
		&model.NewFollow{Following: model.User{Uuid: followUuid.String()}},
	)
}

func (t *TestService) GetPostsForUserFollows(username string, viewerUserUuid uuid.UUID, limit int) ([]*model.Post, error) {
	return t.postService.GetPostsForUserFollows(username, viewerUserUuid, limit)
}

func (t *TestService) GetUserFollows(username string) ([]*model.Follow, error) {
	return t.followService.GetUserFollows(username)
}

func (t *TestService) GetUserFollowers(username string) ([]*model.Follow, error) {
	return t.followService.GetUserFollowers(username)
}
