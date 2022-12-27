package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/util"
	"net/http"
)

// GetShareV1 - get a share
func GetShareV1(c *gin.Context) {
	uuidParam, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.Status(http.StatusOK)
		return
	}
	share, err := service.CreateShareService().GetShare(uuidParam)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, share)
}

// GetSharesV1 - get shares
func GetSharesV1(c *gin.Context) {
	session, _ := util.GetSession(c)
	var viewerUsername string
	if session != nil {
		viewerUser, _ := service.CreateUserService().GetUser(uuid.MustParse(session.User.Uuid))
		viewerUsername = viewerUser.Username
	}
	limit := constants.UserPostsDefaultPageSize
	share, err := service.CreateShareService().GetShares(&viewerUsername, limit)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, share)
}

// CreateShareV1 - create a share
func CreateShareV1(c *gin.Context) {
	newShareParam := model.DecodeRequestToNewShare(c.Request)
	share, err := service.CreateShareService().CreateShare(newShareParam)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusCreated, share)
}
