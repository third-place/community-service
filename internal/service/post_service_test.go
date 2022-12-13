package service

import (
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/util"
	"testing"
)

const message = "this is a test"

func Test_PostService_CreatePublic_NewPost(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)

	// when
	response, err := svc.CreatePost(session, model.CreateNewPost(message))

	// then
	util.Assert(t, err == nil)
	util.Assert(t, response != nil)
}

func Test_PostService_Respects_ProtectedVisibility(t *testing.T) {
	// setup
	svc := CreateTestService()
	user1Model := util.CreateTestUser()

	// given user 1 is protected
	user1Model.Visibility = model.PROTECTED

	// and a few more users...
	user1 := svc.CreateUser(user1Model)
	session1 := model.CreateSessionModelFromString(*user1.Uuid)
	user2 := svc.CreateUser(util.CreateTestUser())
	session2 := model.CreateSessionModelFromString(*user2.Uuid)
	user3 := svc.CreateUser(util.CreateTestUser())
	session3 := model.CreateSessionModelFromString(*user3.Uuid)

	// and user 1 follows user 2
	_, _ = svc.CreateFollow(*user1.Uuid, *user2.Uuid)

	// and user 1 creates a post
	response, _ := svc.CreatePost(session1, model.CreateNewPost(message))

	// when user 2 and 3 get user 1's posts
	post1, err1 := svc.GetPost(session2, uuid.MustParse(response.Uuid))
	post2, err2 := svc.GetPost(session3, uuid.MustParse(response.Uuid))

	// then user 2 can see
	util.Assert(t, post1 != nil)
	util.Assert(t, err1 == nil)

	// and user 3 cannot
	util.Assert(t, post2 == nil)
	util.Assert(t, err2 != nil)
}

func Test_PostService_CreateNewPost_Fails_WhenUserNotFound(t *testing.T) {
	// setup
	svc := CreateTestService()
	userUuid, _ := uuid.NewRandom()
	session := model.CreateSessionModelFromString(userUuid)

	// when
	response, err := svc.CreatePost(session, model.CreateNewPost(message))

	// then
	util.Assert(t, err != nil)
	util.Assert(t, response == nil)
}

func Test_PostService_Can_DeletePost(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)
	postModel, _ := svc.CreatePost(session, model.CreateNewPost(message))

	// when
	err := svc.DeletePost(session, uuid.MustParse(postModel.Uuid))

	// then
	util.Assert(t, err == nil)
}

func Test_PostService_CannotGet_DeletedPost(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)
	postModel, _ := svc.CreatePost(session, model.CreateNewPost(message))
	_ = svc.DeletePost(session, uuid.MustParse(postModel.Uuid))

	// when
	response, err := svc.GetPost(nil, uuid.MustParse(postModel.Uuid))

	// then
	util.Assert(t, err != nil)
	util.Assert(t, response == nil)
}

func Test_GetPosts(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)

	// when
	posts, err := svc.GetPostsFirehose(session, constants.UserPostsDefaultPageSize)

	// then
	util.Assert(t, err == nil)
	util.Assert(t, posts != nil)
}

func Test_GetPost(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)

	// given
	post, err := svc.CreatePost(session, model.CreateNewPost(message))

	// expect
	util.Assert(t, post != nil)
	util.Assert(t, err == nil)

	// when
	response, err := svc.GetPost(nil, uuid.MustParse(post.Uuid))

	// then
	util.Assert(t, err == nil)
	util.Assert(t, response != nil && response.Text == message)
}

func Test_GetPost_Fails_WhenNotFound(t *testing.T) {
	// setup
	svc := CreateTestService()

	// when
	post, err := svc.GetPost(nil, uuid.New())

	// then
	util.Assert(t, err != nil)
	util.Assert(t, post == nil)
}

func Test_PostService_GetUserPosts(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)

	// given
	for i := 0; i < 5; i++ {
		_, _ = svc.CreatePost(session, model.CreateNewPost(message))
	}

	// when
	posts, _ := svc.GetPostsForUser(session, user.Username, constants.UserPostsDefaultPageSize)

	// then
	util.Assert(t, len(posts) == 5)
}

func Test_PostService_GetUserPosts_FailsFor_MissingUser(t *testing.T) {
	// setup
	svc := CreateTestService()
	testUserUuid, _ := uuid.NewRandom()
	session := model.CreateSessionModelFromString(testUserUuid)

	// given
	for i := 0; i < 5; i++ {
		_, _ = svc.CreatePost(
			session,
			model.CreateNewPost(message),
		)
	}

	// when
	posts, err := svc.GetPostsForUser(session, testUserUuid.String(), constants.UserPostsDefaultPageSize)

	// then
	util.Assert(t, posts == nil)
	util.Assert(t, err != nil)
}

func Test_CanGetPosts_ForUserFollows(t *testing.T) {
	// setup
	svc := CreateTestService()
	bob := svc.CreateUser(util.CreateTestUser())
	alice := svc.CreateUser(util.CreateTestUser())

	// given -- bob follows alice
	_, _ = svc.CreateFollow(*bob.Uuid, *alice.Uuid)

	// given -- alice creates some posts
	session := model.CreateSessionModelFromString(*alice.Uuid)
	for i := 0; i < 5; i++ {
		_, _ = svc.CreatePost(session, model.CreateNewPost(message))
	}

	// when -- bob gets posts from people he follows
	posts, err := svc.GetPostsForUserFollows(
		model.CreateSessionModelFromString(*bob.Uuid),
		bob.Username,
		constants.UserPostsDefaultPageSize,
	)

	// then -- expect to see posts from alice
	util.Assert(t, err == nil)
	util.Assert(t, len(posts) >= 5)
}
