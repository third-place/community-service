package service

import (
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/util"
	"testing"
)

func Test_PostReport_HappyPath(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)
	post, _ := svc.CreatePost(session, model.CreateNewPost(""))
	user2 := svc.CreateUser(util.CreateTestUser())

	// when
	report, err := svc.CreatePostReport(
		model.CreateNewPostReport(
			*user2.Uuid,
			uuid.MustParse(post.Uuid),
			"this is offensive",
		),
	)

	// then
	util.Assert(t, err == nil)
	util.Assert(t, report != nil)
}

func Test_PostReport_Fails_WhenPostMissing(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())

	// when
	report, err := svc.CreatePostReport(
		model.CreateNewPostReport(
			*user.Uuid,
			uuid.New(),
			"this is offensive",
		),
	)

	// then
	util.Assert(t, err != nil)
	util.Assert(t, err.Error() == constants.ErrorMessagePostNotFound)
	util.Assert(t, report == nil)
}

func Test_PostReport_Fails_WhenUserMissing(t *testing.T) {
	// setup
	svc := CreateTestService()

	// when
	report, err := svc.CreatePostReport(
		model.CreateNewPostReport(
			uuid.New(),
			uuid.New(),
			"this is offensive",
		),
	)

	// then
	util.Assert(t, err != nil)
	util.Assert(t, err.Error() == constants.ErrorMessageUserNotFound)
	util.Assert(t, report == nil)
}

func Test_ReplyReport_HappyPath(t *testing.T) {
	// setup
	svc := CreateTestService()
	user1 := svc.CreateUser(util.CreateTestUser())
	session1 := model.CreateSessionModelFromString(*user1.Uuid)
	post, _ := svc.CreatePost(session1, model.CreateNewPost(""))
	reply, _ := svc.CreateReply(session1, model.CreateNewReply(uuid.MustParse(post.Uuid), "test message"))
	user2 := svc.CreateUser(util.CreateTestUser())

	// when
	report, err := svc.CreateReplyReport(
		model.CreateNewPostReport(
			*user2.Uuid,
			uuid.MustParse(reply.Uuid),
			"this is offensive",
		),
	)

	// then
	util.Assert(t, err == nil)
	util.Assert(t, report != nil)
}
