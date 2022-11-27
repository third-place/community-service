package service_test

import (
	"github.com/google/uuid"
	model2 "github.com/third-place/community-service/internal/auth/model"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/test"
	"testing"
)

func Test_PostReport_HappyPath(t *testing.T) {
	// setup
	user1 := createTestUser()
	session := model2.CreateSessionModelFromString(*user1.Uuid)
	postService := service.CreatePostService()
	post, _ := postService.CreatePost(session, model.CreateNewPost(""))
	user2 := createTestUser()
	reportService := service.CreateDefaultReportService()
	postUuid := uuid.MustParse(post.Uuid)

	// when
	report, err := reportService.CreatePostReport(model.CreateNewPostReport(user2.Uuid, &postUuid, "this is offensive"))

	// then
	test.Assert(t, err == nil)
	test.Assert(t, report != nil)
}

func Test_PostReport_Fails_WhenPostMissing(t *testing.T) {
	// setup
	postUuid := uuid.New()
	user := createTestUser()
	reportService := service.CreateDefaultReportService()

	// when
	report, err := reportService.CreatePostReport(model.CreateNewPostReport(user.Uuid, &postUuid, "this is offensive"))

	// then
	test.Assert(t, err != nil)
	test.Assert(t, err.Error() == constants.ErrorMessagePostNotFound)
	test.Assert(t, report == nil)
}

func Test_PostReport_Fails_WhenUserMissing(t *testing.T) {
	// setup
	postUuid := uuid.New()
	userUuid := uuid.New()
	reportService := service.CreateDefaultReportService()

	// when
	report, err := reportService.CreatePostReport(model.CreateNewPostReport(&userUuid, &postUuid, "this is offensive"))

	// then
	test.Assert(t, err != nil)
	test.Assert(t, err.Error() == constants.ErrorMessageUserNotFound)
	test.Assert(t, report == nil)
}

func Test_ReplyReport_HappyPath(t *testing.T) {
	// setup
	user1 := createTestUser()
	session1 := model2.CreateSessionModelFromString(*user1.Uuid)
	postService := service.CreatePostService()
	post, _ := postService.CreatePost(session1, model.CreateNewPost(""))
	replyService := service.CreateReplyService()
	postUuid := uuid.MustParse(post.Uuid)
	reply, _ := replyService.CreateReply(session1, model.CreateNewReply(&postUuid, "test message"))
	user2 := createTestUser()
	reportService := service.CreateDefaultReportService()
	replyUuid := uuid.MustParse(reply.Uuid)

	// when
	report, err := reportService.CreateReplyReport(model.CreateNewPostReport(user2.Uuid, &replyUuid, "this is offensive"))

	// then
	test.Assert(t, err == nil)
	test.Assert(t, report != nil)
}
