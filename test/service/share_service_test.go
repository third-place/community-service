package service_test

import (
	model2 "github.com/third-place/community-service/internal/auth/model"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/test"
	"testing"
)

func Test_ShareService_CanCreate_NewShare(t *testing.T) {
	// setup
	testUser := createTestUser()
	session := model2.CreateSessionModelFromString(*testUser.Uuid)
	postService := service.CreatePostService()
	shareService := service.CreateDefaultShareService()
	post, _ := postService.CreatePost(session, model.CreateNewPost(message))
	newShare := &model.NewShare{
		Text: "Yo",
		User: model.User{
			Uuid: testUser.Uuid.String(),
		},
		Post: *post,
	}

	// when
	share, err := shareService.CreateShare(newShare)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, share != nil)
}
