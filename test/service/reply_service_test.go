package service_test

import (
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/entity"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/test"
	"testing"
)

const NumberOfRepliesToCreate = 5

func createTestUser() *entity.User {
	return service.CreateUserService().CreateUser(test.CreateTestUser())
}

func createReplyModel(post *model.Post) *model.NewReply {
	return &model.NewReply{
		Post: *post,
		Text: "this is a reply",
	}
}

func Test_GetReplies_ForPost(t *testing.T) {
	// setup
	testUser := createTestUser()
	session := model.CreateSessionModelFromString(*testUser.Uuid)
	postService := service.CreatePostService()
	replyService := service.CreateReplyService()
	post, err := postService.CreatePost(session, model.CreateNewPost("this is a test"))

	// expect
	test.Assert(t, err == nil)
	test.Assert(t, post != nil)

	// given
	for i := 0; i < NumberOfRepliesToCreate; i++ {
		_, _ = replyService.CreateReply(session, createReplyModel(post))
	}

	// when
	replies, _ := replyService.GetRepliesForPost(uuid.MustParse(post.Uuid))

	// then
	test.Assert(t, len(replies) == NumberOfRepliesToCreate)
}

func Test_CreateReply_Fails_WithMissing_User(t *testing.T) {
	// setup
	testUser, _ := uuid.NewRandom()
	session := model.CreateSessionModelFromString(testUser)
	replyService := service.CreateReplyService()

	// when
	postUuid := uuid.New()
	response, err := replyService.CreateReply(session, &model.NewReply{
		Post: model.Post{
			Uuid: postUuid.String(),
			Text: "",
		},
		Text: "this is a reply",
	})

	// then
	test.Assert(t, err != nil)
	test.Assert(t, err.Error() == constants.ErrorMessageUserNotFound)
	test.Assert(t, response == nil)
}

func Test_CreateReply_Fails_WithMissing_Post(t *testing.T) {
	// setup
	testUser := createTestUser()
	session := model.CreateSessionModelFromString(*testUser.Uuid)
	replyService := service.CreateReplyService()

	// when
	postUuid := uuid.New()
	response, err := replyService.CreateReply(session, &model.NewReply{
		Post: model.Post{
			Uuid: postUuid.String(),
			Text: "",
		},
		Text: "this is a reply",
	})

	// then
	test.Assert(t, err != nil)
	test.Assert(t, err.Error() == constants.ErrorMessagePostNotFound)
	test.Assert(t, response == nil)
}

func Test_GetReplies_FailsWithMissing_Post(t *testing.T) {
	// setup
	replyService := service.CreateReplyService()

	// when
	response, err := replyService.GetRepliesForPost(uuid.New())

	// then
	test.Assert(t, err != nil)
	test.Assert(t, err.Error() == constants.ErrorMessagePostNotFound)
	test.Assert(t, response == nil)
}
