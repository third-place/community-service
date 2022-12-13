package service

import (
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/entity"
	"github.com/third-place/community-service/internal/model"
)

type TestService struct {
	userService   *UserService
	postService   *PostService
	replyService  *ReplyService
	followService *FollowService
	reportService *ReportService
	ShareService  *ShareService
}

func CreateTestService() *TestService {
	return &TestService{
		CreateTestUserService(),
		CreateTestPostService(),
		CreateTestReplyService(),
		CreateTestFollowService(),
		CreateTestReportService(),
		CreateTestShareService(),
	}
}

func (t *TestService) CreateUser(user *model.User) *entity.User {
	return t.userService.CreateUser(user)
}

func (t *TestService) GetUser(userUuid uuid.UUID) (*model.User, error) {
	return t.userService.GetUser(userUuid)
}

func (t *TestService) DeleteUser(userUuid uuid.UUID) error {
	return t.userService.DeleteUser(userUuid)
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

func (t *TestService) GetPostsFirehose(session *model.Session, limit int) ([]*model.Post, error) {
	return t.postService.GetPostsFirehose(session, limit)
}

func (t *TestService) GetPostsForUser(session *model.Session, username string, limit int) ([]*model.Post, error) {
	return t.postService.GetPostsForUser(session, username, limit)
}

func (t *TestService) CreateFollow(sessionUserUuid uuid.UUID, followUuid uuid.UUID) (*model.Follow, error) {
	return t.followService.CreateFollow(
		sessionUserUuid,
		&model.NewFollow{Following: model.User{Uuid: followUuid.String()}},
	)
}

func (t *TestService) GetPostsForUserFollows(session *model.Session, username string, limit int) ([]*model.Post, error) {
	return t.postService.GetPostsForUserFollows(session, username, limit)
}

func (t *TestService) GetUserFollows(username string) ([]*model.Follow, error) {
	return t.followService.GetUserFollows(username)
}

func (t *TestService) GetUserFollowers(username string) ([]*model.Follow, error) {
	return t.followService.GetUserFollowers(username)
}

func (t *TestService) CreateReply(session *model.Session, reply *model.NewReply) (*model.Post, error) {
	return t.replyService.CreateReply(session, reply)
}

func (t *TestService) GetRepliesForPost(postUuid uuid.UUID) ([]*model.Post, error) {
	return t.replyService.GetRepliesForPost(postUuid)
}

func (t *TestService) CreatePostReport(newReport *model.NewPostReport) (*model.PostReport, error) {
	return t.reportService.CreatePostReport(newReport)
}

func (t *TestService) CreateReplyReport(newReport *model.NewPostReport) (*model.PostReport, error) {
	return t.reportService.CreateReplyReport(newReport)
}

func (t *TestService) CreateShare(share *model.NewShare) (*model.Share, error) {
	return t.ShareService.CreateShare(share)
}
