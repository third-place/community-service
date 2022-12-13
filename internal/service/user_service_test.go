package service

import (
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/util"
	"testing"
)

func Test_CanGetUser(t *testing.T) {
	// setup
	svc := CreateTestService()
	user := svc.CreateUser(util.CreateTestUser())

	// when
	response, err := svc.GetUser(*user.Uuid)

	// then
	util.Assert(t, err == nil)
	util.Assert(t, response.Uuid == user.Uuid.String())
}

func Test_DeleteMissingUser_Fails(t *testing.T) {
	// setup
	svc := CreateTestService()

	// when
	err := svc.DeleteUser(uuid.New())

	// then
	util.Assert(t, err != nil)
	util.Assert(t, err.Error() == constants.ErrorMessageUserNotFound)
}

func Test_Requesting_UserNotFound(t *testing.T) {
	// setup
	svc := CreateTestService()

	// when
	user, err := svc.GetUser(uuid.New())

	// then
	util.Assert(t, user == nil)
	util.Assert(t, err != nil)
}
