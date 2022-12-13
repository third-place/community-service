package service

import (
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/test"
	"github.com/third-place/community-service/internal/util"
	"testing"
)

const NumberOfRepliesToCreate = 5

func createReplyModel(post *model.Post) *model.NewReply {
	return &model.NewReply{
		Post: *post,
		Text: "this is a reply",
	}
}

func Test_GetReplies_ForPost(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)
	post, err := svc.CreatePost(session, model.CreateNewPost("this is a test"))

	// expect
	test.Assert(t, err == nil)
	test.Assert(t, post != nil)

	// given
	for i := 0; i < NumberOfRepliesToCreate; i++ {
		_, _ = svc.CreateReply(session, createReplyModel(post))
	}

	// when
	replies, _ := svc.GetRepliesForPost(uuid.MustParse(post.Uuid))

	// then
	test.Assert(t, len(replies) == NumberOfRepliesToCreate)
}

func Test_CreateReply_Fails_WithMissing_User(t *testing.T) {
	// setup
	svc := CreateTestService()
	testUser, _ := uuid.NewRandom()
	session := model.CreateSessionModelFromString(testUser)

	// when
	postUuid := uuid.New()
	response, err := svc.CreateReply(session, &model.NewReply{
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
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)

	// when
	postUuid := uuid.New()
	response, err := svc.CreateReply(session, &model.NewReply{
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
	svc := CreateTestService()

	// when
	response, err := svc.GetRepliesForPost(uuid.New())

	// then
	test.Assert(t, err != nil)
	test.Assert(t, err.Error() == constants.ErrorMessagePostNotFound)
	test.Assert(t, response == nil)
}
