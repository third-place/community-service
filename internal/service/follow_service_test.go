package service

import (
	"github.com/third-place/community-service/internal/test"
	"testing"
)

func Test_CanFollow_User(t *testing.T) {
	// setup
	svc := CreateTestService()
	user1 := svc.CreateUser(test.CreateTestUser())
	user2 := svc.CreateUser(test.CreateTestUser())

	// when
	follow, err := svc.CreateFollow(*user1.Uuid, *user2.Uuid)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, follow.User.Uuid == user1.Uuid.String())
	test.Assert(t, follow.Following.Uuid == user2.Uuid.String())
}

func Test_GetFollows(t *testing.T) {
	// setup
	svc := CreateTestService()
	user1 := svc.CreateUser(test.CreateTestUser())
	user2 := svc.CreateUser(test.CreateTestUser())
	user3 := svc.CreateUser(test.CreateTestUser())

	_, _ = svc.CreateFollow(*user1.Uuid, *user2.Uuid)
	_, _ = svc.CreateFollow(*user1.Uuid, *user3.Uuid)

	following, err := svc.GetUserFollows(user1.Username)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, len(following) == 2)
}

func Test_GetFollowers(t *testing.T) {
	// setup
	svc := CreateTestService()
	user1 := svc.CreateUser(test.CreateTestUser())
	user2 := svc.CreateUser(test.CreateTestUser())
	user3 := svc.CreateUser(test.CreateTestUser())

	_, _ = svc.CreateFollow(*user1.Uuid, *user3.Uuid)
	_, _ = svc.CreateFollow(*user2.Uuid, *user3.Uuid)

	followers, err := svc.GetUserFollowers(user3.Username)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, len(followers) == 2)
}
