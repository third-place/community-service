package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/util"
	"net/http"
)

// CreateReplyV1 - create a reply
func CreateReplyV1(c *gin.Context) {
	session, err := util.GetSession(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	newReplyModel := model.DecodeRequestToNewReply(c.Request)
	reply, err := service.CreateReplyService().CreateReply(session, newReplyModel)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
	c.JSON(http.StatusOK, reply)
}

func GetPostRepliesV1(c *gin.Context) {
	postUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	replies, err := service.CreateReplyService().GetRepliesForPost(postUuid)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, replies)
}
