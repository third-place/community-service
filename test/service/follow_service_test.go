package service

import (
	"github.com/danielmunro/otto-community-service/internal/model"
	"github.com/danielmunro/otto-community-service/internal/service"
	"github.com/danielmunro/otto-community-service/internal/test"
	"github.com/google/uuid"
	"testing"
)

func Test_CanFollow_AUser(t *testing.T) {
	// setup
	userService := service.CreateDefaultUserService()
	user1 := userService.CreateUser(test.CreateTestUser())
	user2 := userService.CreateUser(test.CreateTestUser())
	followService := service.CreateFollowService()

	// when
	follow, err := followService.CreateFollow(uuid.MustParse(user1.Uuid.String()), &model.NewFollow{
		Following: model.User{Uuid: user2.Uuid.String()},
	})

	// then
	test.Assert(t, err == nil)
	test.Assert(t, follow.User.Uuid == user1.Uuid.String())
	test.Assert(t, follow.Following.Uuid == user2.Uuid.String())
}

func Test_GetFollows(t *testing.T) {
	// setup
	userService := service.CreateDefaultUserService()
	user1 := userService.CreateUser(test.CreateTestUser())
	user2 := userService.CreateUser(test.CreateTestUser())
	user3 := userService.CreateUser(test.CreateTestUser())
	followService := service.CreateFollowService()

	_, _ = followService.CreateFollow(uuid.MustParse(user1.Uuid.String()), &model.NewFollow{
		Following: model.User{Uuid: user2.Uuid.String()},
	})
	_, _ = followService.CreateFollow(uuid.MustParse(user1.Uuid.String()), &model.NewFollow{
		Following: model.User{Uuid: user3.Uuid.String()},
	})

	following, err := followService.GetUserFollows(user1.Username)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, len(following) == 2)
}

func Test_GetFollowers(t *testing.T) {
	// setup
	userService := service.CreateDefaultUserService()
	user1 := userService.CreateUser(test.CreateTestUser())
	user2 := userService.CreateUser(test.CreateTestUser())
	user3 := userService.CreateUser(test.CreateTestUser())
	followService := service.CreateFollowService()

	_, _ = followService.CreateFollow(uuid.MustParse(user1.Uuid.String()), &model.NewFollow{
		Following: model.User{Uuid: user3.Uuid.String()},
	})
	_, _ = followService.CreateFollow(uuid.MustParse(user2.Uuid.String()), &model.NewFollow{
		Following: model.User{Uuid: user3.Uuid.String()},
	})

	followers, err := followService.GetUserFollowers(user3.Username)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, len(followers) == 2)
}
