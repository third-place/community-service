package service

import (
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/test"
	"github.com/third-place/community-service/internal/util"
	"testing"
)

func Test_ShareService_CanCreate_NewShare(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)
	post, _ := svc.CreatePost(session, model.CreateNewPost(message))
	newShare := &model.NewShare{
		Text: "Yo",
		User: model.User{
			Uuid: user.Uuid.String(),
		},
		Post: *post,
	}

	// when
	share, err := svc.CreateShare(newShare)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, share != nil)
}
