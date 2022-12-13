package service

import (
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/util"
	"testing"
)

func Test_ShareService_CanCreate_NewShare(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSession(*user.Uuid)
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
	util.Assert(t, err == nil)
	util.Assert(t, share != nil)
}
