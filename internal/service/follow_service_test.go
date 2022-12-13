package service

import (
	"github.com/third-place/community-service/internal/util"
	"testing"
)

func Test_CanFollow_User(t *testing.T) {
	// setup
	svc := CreateTestService()
	user1 := svc.CreateUser(util.CreateTestUser())
	user2 := svc.CreateUser(util.CreateTestUser())

	// when
	follow, err := svc.CreateFollow(*user1.Uuid, *user2.Uuid)

	// then
	util.Assert(t, err == nil)
	util.Assert(t, follow.User.Uuid == user1.Uuid.String())
	util.Assert(t, follow.Following.Uuid == user2.Uuid.String())
}

func Test_GetFollows(t *testing.T) {
	// setup
	svc := CreateTestService()
	user1 := svc.CreateUser(util.CreateTestUser())
	user2 := svc.CreateUser(util.CreateTestUser())
	user3 := svc.CreateUser(util.CreateTestUser())

	_, _ = svc.CreateFollow(*user1.Uuid, *user2.Uuid)
	_, _ = svc.CreateFollow(*user1.Uuid, *user3.Uuid)

	following, err := svc.GetUserFollows(user1.Username)

	// then
	util.Assert(t, err == nil)
	util.Assert(t, len(following) == 2)
}

func Test_GetFollowers(t *testing.T) {
	// setup
	svc := CreateTestService()
	user1 := svc.CreateUser(util.CreateTestUser())
	user2 := svc.CreateUser(util.CreateTestUser())
	user3 := svc.CreateUser(util.CreateTestUser())

	_, _ = svc.CreateFollow(*user1.Uuid, *user3.Uuid)
	_, _ = svc.CreateFollow(*user2.Uuid, *user3.Uuid)

	followers, err := svc.GetUserFollowers(user3.Username)

	// then
	util.Assert(t, err == nil)
	util.Assert(t, len(followers) == 2)
}
