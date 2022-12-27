package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/util"
	"net/http"
)

// CreateNewPostLikeV1 - like a post
func CreateNewPostLikeV1(c *gin.Context) {
	uuidParam, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	session, err := util.GetSession(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	postLike, err := service.CreateDefaultLikeService().CreateLikeForPost(uuidParam, uuid.MustParse(session.User.Uuid))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusCreated, postLike)
}

// DeleteLikeForPostV1 - delete like for post
func DeleteLikeForPostV1(c *gin.Context) {
	uuidParam, err := uuid.Parse(c.Param("uuid"))
	session, err := util.GetSession(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	err = service.CreateDefaultLikeService().DeleteLikeForPost(uuidParam, uuid.MustParse(session.User.Uuid))
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
}
