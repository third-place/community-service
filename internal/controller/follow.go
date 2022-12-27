package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/mapper"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/util"
	"log"
	"net/http"
)

// CreateFollowV1 - create a follow
func CreateFollowV1(c *gin.Context) {
	newFollowModel, err := mapper.DecodeRequestToNewFollow(c.Request)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	session, err := util.GetSession(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	follow, err := service.CreateFollowService().CreateFollow(uuid.MustParse(session.User.Uuid), newFollowModel)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusCreated, follow)
}

// GetUserFollowersV1 - get user followers
func GetUserFollowersV1(c *gin.Context) {
	username := c.Param("username")
	follows, err := service.CreateFollowService().GetUserFollowers(username)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, follows)
}

// GetUserFollowsV1 - get user follows
func GetUserFollowsV1(c *gin.Context) {
	username := c.Param("username")
	follows, err := service.CreateFollowService().GetUserFollows(username)
	if err != nil {
		log.Print("error received from get user follows :: ", err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, follows)
}

// DeleteFollowV1 - delete a follow
func DeleteFollowV1(c *gin.Context) {
	session, err := util.GetSession(c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	followUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = service.CreateFollowService().DeleteFollow(followUuid, uuid.MustParse(session.User.Uuid))
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
}
