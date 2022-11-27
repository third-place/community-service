package service

import (
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/test"
	"github.com/google/uuid"
	"testing"
)

func Test_CanGetUser(t *testing.T) {
	// setup
	userService := service.CreateDefaultUserService()
	user := userService.CreateUser(test.CreateTestUser())

	// when
	response, err := userService.GetUser(*user.Uuid)

	// then
	test.Assert(t, err == nil)
	test.Assert(t, response.Uuid == user.Uuid.String())
}

func Test_DeleteMissingUser_Fails(t *testing.T) {
	// setup
	userService := service.CreateDefaultUserService()
	userModel := test.CreateTestUser()

	// when
	err := userService.DeleteUser(uuid.MustParse(userModel.Uuid))

	// then
	test.Assert(t, err != nil)
	test.Assert(t, err.Error() == constants.ErrorMessageUserNotFound)
}

func Test_Requesting_UserNotFound(t *testing.T) {
	// setup
	userService := service.CreateDefaultUserService()

	// when
	user, err := userService.GetUser(uuid.New())

	// then
	test.Assert(t, user == nil)
	test.Assert(t, err != nil)
}
