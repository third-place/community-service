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
	session := svc.CreateTestUserSession()

	// when
	response, err := svc.CreatePost(session, model.CreateNewPost(message))

	// then
	if err != nil {
		t.Error(err)
	}
	if response == nil {
		t.Fail()
	}
}

func Test_BannedUser_CannotCreatePost(t *testing.T) {
	// setup
	svc := CreateTestService()
	session := svc.CreateTestBannedUserSession()

	// when
	response, err := svc.CreatePost(session, model.CreateNewPost(message))

	// then
	if response != nil || err == nil || err.Error() != "not allowed to create a post" {
		t.Fail()
	}
}

func Test_PostService_Respects_ProtectedVisibility(t *testing.T) {
	// setup
	svc := CreateTestService()
	user1Model := util.CreateTestUser()

	// given user 1 is protected
	user1Model.Visibility = model.PROTECTED

	// and a few more users...
	user1 := svc.CreateUser(user1Model)
	session1 := model.CreateSession(*user1.Uuid)
	user2 := svc.CreateUser(util.CreateTestUser())
	session2 := model.CreateSession(*user2.Uuid)
	user3 := svc.CreateUser(util.CreateTestUser())
	session3 := model.CreateSession(*user3.Uuid)

	// and user 1 follows user 2
	_, _ = svc.CreateFollow(*user1.Uuid, *user2.Uuid)

	// and user 1 creates a post
	response, _ := svc.CreatePost(session1, model.CreateNewPost(message))

	// when user 2 and 3 get user 1's posts
	post1, err1 := svc.GetPost(session2, uuid.MustParse(response.Uuid))
	post2, err2 := svc.GetPost(session3, uuid.MustParse(response.Uuid))

	// then user 2 can see
	if post1 == nil || err1 != nil {
		t.Fail()
	}

	// and user 3 cannot
	if post2 != nil || err2 == nil {
		t.Fail()
	}
}

func Test_PostService_Respects_PrivateVisibility(t *testing.T) {
	// setup
	svc := CreateTestService()
	user1Model := util.CreateTestUser()

	// given user 1 is private
	user1Model.Visibility = model.PRIVATE

	// and a few more users...
	user1 := svc.CreateUser(user1Model)
	session1 := model.CreateSession(*user1.Uuid)
	user2 := svc.CreateUser(util.CreateTestUser())
	session2 := model.CreateSession(*user2.Uuid)
	user3 := svc.CreateUser(util.CreateTestUser())
	session3 := model.CreateSession(*user3.Uuid)

	// and user 1 follows user 2
	_, _ = svc.CreateFollow(*user1.Uuid, *user2.Uuid)

	// and user 1 creates a post
	response, _ := svc.CreatePost(session1, model.CreateNewPost(message))

	// when user 2 and 3 get user 1's posts
	post1, err1 := svc.GetPost(session1, uuid.MustParse(response.Uuid))
	post2, err2 := svc.GetPost(session2, uuid.MustParse(response.Uuid))
	post3, err3 := svc.GetPost(session3, uuid.MustParse(response.Uuid))

	// then user 1 can see
	if post1 == nil || err1 != nil {
		t.Fail()
	}

	// and user 2 cannot see
	if post2 != nil || err2 == nil {
		t.Fail()
	}

	// and user 3 cannot see
	if post3 != nil || err3 == nil {
		t.Fail()
	}
}

func Test_PostService_CreateNewPost_Fails_WhenUserNotFound(t *testing.T) {
	// setup
	svc := CreateTestService()
	userUuid, _ := uuid.NewRandom()
	session := model.CreateSession(userUuid)

	// when
	response, err := svc.CreatePost(session, model.CreateNewPost(message))

	// then
	if response != nil || err == nil {
		t.Fail()
	}
}

func Test_PostService_Can_DeletePost(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSession(*user.Uuid)
	postModel, _ := svc.CreatePost(session, model.CreateNewPost(message))

	// when
	err := svc.DeletePost(session, uuid.MustParse(postModel.Uuid))

	// then
	if err != nil {
		t.Error(err)
	}
}

func Test_PostService_CannotGet_DeletedPost(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSession(*user.Uuid)
	postModel, _ := svc.CreatePost(session, model.CreateNewPost(message))
	_ = svc.DeletePost(session, uuid.MustParse(postModel.Uuid))

	// when
	response, err := svc.GetPost(nil, uuid.MustParse(postModel.Uuid))

	// then
	if response != nil || err == nil {
		t.Fail()
	}
}

func Test_GetPosts(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSession(*user.Uuid)

	// when
	posts, err := svc.GetPostsFirehose(session, constants.UserPostsDefaultPageSize)

	// then
	if posts == nil || err != nil {
		t.Fail()
	}
}

func Test_GetPost(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSession(*user.Uuid)

	// given
	post, err := svc.CreatePost(session, model.CreateNewPost(message))

	// expect
	if post == nil || err != nil {
		t.Fail()
	}

	// when
	response, err := svc.GetPost(nil, uuid.MustParse(post.Uuid))

	// then
	if err != nil {
		t.Error(err)
	}
	if response == nil || response.Text != message {
		t.Fail()
	}
}

func Test_GetPost_Fails_WhenNotFound(t *testing.T) {
	// setup
	svc := CreateTestService()

	// when
	post, err := svc.GetPost(nil, uuid.New())

	// then
	if post != nil || err == nil {
		t.Fail()
	}
}

func Test_PostService_GetUserPosts(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())
	session := model.CreateSession(*user.Uuid)

	// given
	for i := 0; i < 5; i++ {
		_, _ = svc.CreatePost(session, model.CreateNewPost(message))
	}

	// when
	posts, _ := svc.GetPostsForUser(session, user.Username, constants.UserPostsDefaultPageSize)

	// then
	if len(posts) != 5 {
		t.Fail()
	}
}

func Test_PostService_GetUserPosts_FailsFor_MissingUser(t *testing.T) {
	// setup
	svc := CreateTestService()
	testUserUuid, _ := uuid.NewRandom()
	session := model.CreateSession(testUserUuid)

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
	if posts != nil || err == nil {
		t.Fail()
	}
}

func Test_CanGetPosts_ForUserFollows(t *testing.T) {
	// setup
	svc := CreateTestService()
	bob := svc.CreateUser(util.CreateTestUser())
	alice := svc.CreateUser(util.CreateTestUser())

	// given -- bob follows alice
	_, _ = svc.CreateFollow(*bob.Uuid, *alice.Uuid)

	// given -- alice creates some posts
	session := model.CreateSession(*alice.Uuid)
	for i := 0; i < 5; i++ {
		_, _ = svc.CreatePost(session, model.CreateNewPost(message))
	}

	// when -- bob gets posts from people he follows
	posts, err := svc.GetPostsForUserFollows(
		model.CreateSession(*bob.Uuid),
		bob.Username,
		constants.UserPostsDefaultPageSize,
	)

	// then -- expect to see posts from alice
	if err != nil || len(posts) < 5 {
		t.Fail()
	}
}
