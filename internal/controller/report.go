package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/util"
	"net/http"
)

// CreatePostReportV1 - report a post
func CreatePostReportV1(c *gin.Context) {
	_, err := util.GetSession(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	newReport := model.DecodeRequestToNewPostReport(c.Request)
	report, err := service.CreateReportService().CreatePostReport(newReport)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, report)
}

// CreateReplyReportV1 - report a reply
func CreateReplyReportV1(c *gin.Context) {
	_, err := util.GetSession(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	newReport := model.DecodeRequestToNewPostReport(c.Request)
	replyReport, err := service.CreateReportService().CreateReplyReport(newReport)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, replyReport)
}
