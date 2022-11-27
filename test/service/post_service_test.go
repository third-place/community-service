package service_test

import (
	model2 "github.com/third-place/community-service/internal/auth/model"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/entity"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/test"
	"github.com/google/uuid"
	"testing"
)

const message = "this is a test"

func createTestUser() *entity.User {
	return service.CreateDefaultUserService().CreateUser(test.CreateTestUser())
}

func Test_PostService_CreatePublic_NewPost(t *testing.T) {
	// setup
	testUser := createTestUser()
	session := model2.CreateSessionModelFromString(*testUser.Uuid)
	postService := service.CreatePostService()

	// when
	response, err := postService.CreatePost(session, model.CreateNewPost(message))

	// then
	test.Assert(t, err == nil)
	test.Assert(t, response != nil)
	test.Assert(t, response.Visibility == model.PUBLIC)
}

func Test_PostService_CreateNewPost_WithPrivateVisibility(t *testing.T) {
	// setup
	testUser := createTestUser()
	session := model2.CreateSessionModelFromString(*testUser.Uuid)
	postService := service.CreatePostService()

	// given
	newPost := model.CreateNewPost(message)
	newPost.Visibility = model.PRIVATE

	// when
	response, err := postService.CreatePost(session, newPost)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, response.Visibility == model.PRIVATE)
}

func Test_PostService_Respects_PrivateVisibility(t *testing.T) {
	// setup
	testUser := createTestUser()
	session := model2.CreateSessionModelFromString(*testUser.Uuid)
	postService := service.CreatePostService()
	newPost := model.CreateNewPost(message)
	newPost.Visibility = model.PRIVATE
	response, err := postService.CreatePost(session, newPost)

	// when
	post, err := postService.GetPost(nil, uuid.MustParse(response.Uuid))

	// then
	test.Assert(t, post == nil)
	test.Assert(t, err != nil)
}

func Test_PostService_CreateNewPost_WithFollowingVisibility(t *testing.T) {
	// setup
	testUser := createTestUser()
	session := model2.CreateSessionModelFromString(*testUser.Uuid)
	postService := service.CreatePostService()

	// given
	newPost := model.CreateNewPost(message)
	newPost.Visibility = model.FOLLOWING

	// when
	response, err := postService.CreatePost(session, newPost)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, response.Visibility == model.FOLLOWING)
}

func Test_PostService_Respects_FollowingVisibility(t *testing.T) {
	// setup
	testUser1 := createTestUser()
	session1 := model2.CreateSessionModelFromString(*testUser1.Uuid)
	testUser2 := createTestUser()
	session2 := model2.CreateSessionModelFromString(*testUser2.Uuid)
	testUser3 := createTestUser()
	session3 := model2.CreateSessionModelFromString(*testUser3.Uuid)
	_, _ = service.CreateFollowService().CreateFollow(
		*testUser1.Uuid, &model.NewFollow{Following: model.User{Uuid: testUser2.Uuid.String()}})
	postService := service.CreatePostService()
	newPost := model.CreateNewPost(message)
	newPost.Visibility = model.FOLLOWING
	response, _ := postService.CreatePost(session1, newPost)

	// when
	post1, err1 := postService.GetPost(session2, uuid.MustParse(response.Uuid))
	post2, err2 := postService.GetPost(session3, uuid.MustParse(response.Uuid))

	// then
	test.Assert(t, post1 != nil)
	test.Assert(t, err1 == nil)

	test.Assert(t, post2 == nil)
	test.Assert(t, err2 != nil)
}

func Test_PostService_CreateNewPost_Fails_WhenUserNotFound(t *testing.T) {
	// setup
	userUuid, _ := uuid.NewRandom()
	session := model2.CreateSessionModelFromString(userUuid)
	postService := service.CreatePostService()

	// when
	response, err := postService.CreatePost(session, model.CreateNewPost(message))

	// then
	test.Assert(t, err != nil)
	test.Assert(t, response == nil)
}

func Test_PostService_Can_DeletePost(t *testing.T) {
	// setup
	testUser := createTestUser()
	session := model2.CreateSessionModelFromString(*testUser.Uuid)
	postService := service.CreatePostService()
	postModel, _ := postService.CreatePost(session, model.CreateNewPost(message))

	// when
	err := postService.DeletePost(session, uuid.MustParse(postModel.Uuid))

	// then
	test.Assert(t, err == nil)
}

func Test_PostService_CannotGet_DeletedPost(t *testing.T) {
	// setup
	testUser := createTestUser()
	session := model2.CreateSessionModelFromString(*testUser.Uuid)
	postService := service.CreatePostService()
	postModel, _ := postService.CreatePost(session, model.CreateNewPost(message))
	_ = postService.DeletePost(session, uuid.MustParse(postModel.Uuid))

	// when
	response, err := postService.GetPost(nil, uuid.MustParse(postModel.Uuid))

	// then
	test.Assert(t, err != nil)
	test.Assert(t, response == nil)
}

func Test_GetPosts(t *testing.T) {
	// setup
	postService := service.CreatePostService()
	testUser := createTestUser()

	// when
	posts, err := postService.GetPostsFirehose(&testUser.Username, constants.UserPostsDefaultPageSize)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, posts != nil)
}

func Test_GetPosts_NoSession(t *testing.T) {
	// setup
	postService := service.CreatePostService()

	// when
	posts, err := postService.GetPostsFirehose(nil, constants.UserPostsDefaultPageSize)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, posts != nil)
}

func Test_GetPost(t *testing.T) {
	// setup
	testUser := createTestUser()
	session := model2.CreateSessionModelFromString(*testUser.Uuid)
	postService := service.CreatePostService()

	// given
	post, err := postService.CreatePost(session, model.CreateNewPost(message))

	// expect
	test.Assert(t, post != nil)
	test.Assert(t, err == nil)

	// when
	response, err := postService.GetPost(nil, uuid.MustParse(post.Uuid))

	// then
	test.Assert(t, err == nil)
	test.Assert(t, response != nil && response.Text == message)
}

func Test_GetPost_Fails_WhenNotFound(t *testing.T) {
	// setup
	postService := service.CreatePostService()

	// when
	post, err := postService.GetPost(nil, uuid.New())

	// then
	test.Assert(t, err != nil)
	test.Assert(t, post == nil)
}

func Test_PostService_GetUserPosts(t *testing.T) {
	// setup
	testUser := createTestUser()
	session := model2.CreateSessionModelFromString(*testUser.Uuid)
	postService := service.CreatePostService()

	// given
	for i := 0; i < 5; i++ {
		_, _ = postService.CreatePost(session, model.CreateNewPost(message))
	}

	// when
	posts, _ := postService.GetPostsForUser(testUser.Username, nil, constants.UserPostsDefaultPageSize)

	// then
	test.Assert(t, len(posts) == 5)
}

func Test_PostService_GetUserPosts_FailsFor_MissingUser(t *testing.T) {
	// setup
	testUserUuid, _ := uuid.NewRandom()
	session := model2.CreateSessionModelFromString(testUserUuid)
	postService := service.CreatePostService()

	// given
	for i := 0; i < 5; i++ {
		_, _ = postService.CreatePost(
			session,
			model.CreateNewPost(message),
		)
	}

	// when
	posts, err := postService.GetPostsForUser(testUserUuid.String(), nil, constants.UserPostsDefaultPageSize)

	// then
	test.Assert(t, posts == nil)
	test.Assert(t, err != nil)
}

func Test_CanGetPosts_ForUserFollows(t *testing.T) {
	// setup
	testUser := createTestUser()
	following := createTestUser()
	postService := service.CreatePostService()
	followService := service.CreateFollowService()
	_, _ = followService.CreateFollow(
		*testUser.Uuid,
		&model.NewFollow{Following: model.User{Uuid: following.Uuid.String()}},
	)
	session := model2.CreateSessionModelFromString(*following.Uuid)

	// given
	for i := 0; i < 5; i++ {
		_, _ = postService.CreatePost(session, model.CreateNewPost(message))
	}

	// when
	posts, err := postService.GetPostsForUserFollows(testUser.Username, *testUser.Uuid, constants.UserPostsDefaultPageSize)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, len(posts) >= 5)
}
