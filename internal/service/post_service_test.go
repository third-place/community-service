package service

import (
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/test"
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
	test.Assert(t, err == nil)
	test.Assert(t, response != nil)
	test.Assert(t, response.Visibility == model.PUBLIC)
}

func Test_PostService_CreateNewPost_WithPrivateVisibility(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)

	// given
	newPost := model.CreateNewPost(message)
	newPost.Visibility = model.PRIVATE

	// when
	response, err := svc.CreatePost(session, newPost)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, response.Visibility == model.PRIVATE)
}

func Test_PostService_Respects_PrivateVisibility(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)
	newPost := model.CreateNewPost(message)
	newPost.Visibility = model.PRIVATE
	response, err := svc.CreatePost(session, newPost)

	// when
	post, err := svc.GetPost(nil, uuid.MustParse(response.Uuid))

	// then
	test.Assert(t, post == nil)
	test.Assert(t, err != nil)
}

func Test_PostService_CreateNewPost_WithFollowingVisibility(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)

	// given
	newPost := model.CreateNewPost(message)
	newPost.Visibility = model.FOLLOWING

	// when
	response, err := svc.CreatePost(session, newPost)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, response.Visibility == model.FOLLOWING)
}

func Test_PostService_Respects_FollowingVisibility(t *testing.T) {
	// setup
	svc := CreateTestService()
	user1 := svc.CreateUser(util.CreateTestUser())
	session1 := model.CreateSessionModelFromString(*user1.Uuid)
	user2 := svc.CreateUser(util.CreateTestUser())
	session2 := model.CreateSessionModelFromString(*user2.Uuid)
	user3 := svc.CreateUser(util.CreateTestUser())
	session3 := model.CreateSessionModelFromString(*user3.Uuid)
	_, _ = svc.CreateFollow(*user1.Uuid, *user2.Uuid)
	newPost := model.CreateNewPost(message)
	newPost.Visibility = model.FOLLOWING
	response, _ := svc.CreatePost(session1, newPost)

	// when
	post1, err1 := svc.GetPost(session2, uuid.MustParse(response.Uuid))
	post2, err2 := svc.GetPost(session3, uuid.MustParse(response.Uuid))

	// then
	test.Assert(t, post1 != nil)
	test.Assert(t, err1 == nil)

	test.Assert(t, post2 == nil)
	test.Assert(t, err2 != nil)
}

func Test_PostService_CreateNewPost_Fails_WhenUserNotFound(t *testing.T) {
	// setup
	svc := CreateTestService()
	userUuid, _ := uuid.NewRandom()
	session := model.CreateSessionModelFromString(userUuid)

	// when
	response, err := svc.CreatePost(session, model.CreateNewPost(message))

	// then
	test.Assert(t, err != nil)
	test.Assert(t, response == nil)
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
	test.Assert(t, err == nil)
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
	test.Assert(t, err != nil)
	test.Assert(t, response == nil)
}

func Test_GetPosts(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())

	// when
	posts, err := svc.GetPostsFirehose(&user.Username, constants.UserPostsDefaultPageSize)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, posts != nil)
}

func Test_GetPosts_NoSession(t *testing.T) {
	// setup
	svc := CreateTestService()

	// when
	posts, err := svc.GetPostsFirehose(nil, constants.UserPostsDefaultPageSize)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, posts != nil)
}

func Test_GetPost(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSessionModelFromString(*user.Uuid)

	// given
	post, err := svc.CreatePost(session, model.CreateNewPost(message))

	// expect
	test.Assert(t, post != nil)
	test.Assert(t, err == nil)

	// when
	response, err := svc.GetPost(nil, uuid.MustParse(post.Uuid))

	// then
	test.Assert(t, err == nil)
	test.Assert(t, response != nil && response.Text == message)
}

func Test_GetPost_Fails_WhenNotFound(t *testing.T) {
	// setup
	svc := CreateTestService()

	// when
	post, err := svc.GetPost(nil, uuid.New())

	// then
	test.Assert(t, err != nil)
	test.Assert(t, post == nil)
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
	posts, _ := svc.GetPostsForUser(user.Username, nil, constants.UserPostsDefaultPageSize)

	// then
	test.Assert(t, len(posts) == 5)
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
	posts, err := svc.GetPostsForUser(testUserUuid.String(), nil, constants.UserPostsDefaultPageSize)

	// then
	test.Assert(t, posts == nil)
	test.Assert(t, err != nil)
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
	posts, err := svc.GetPostsForUserFollows(bob.Username, *bob.Uuid, constants.UserPostsDefaultPageSize)

	// then -- expect to see posts from alice
	test.Assert(t, err == nil)
	test.Assert(t, len(posts) >= 5)
}
