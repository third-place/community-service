package controller

import (
	"encoding/json"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	"net/http"
)

// CreatePostReportV1 - report a post
func CreatePostReportV1(w http.ResponseWriter, r *http.Request) {
	_, err := service.CreateDefaultAuthService().GetSessionFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	newReport := model.DecodeRequestToNewPostReport(r)
	report, err := service.CreateDefaultReportService().CreatePostReport(newReport)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, _ := json.Marshal(report)
	_, _ = w.Write(data)
}

// CreateReplyReportV1 - report a reply
func CreateReplyReportV1(w http.ResponseWriter, r *http.Request) {
	_, err := service.CreateDefaultAuthService().GetSessionFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	newReport := model.DecodeRequestToNewPostReport(r)
	replyReport, err := service.CreateDefaultReportService().CreateReplyReport(newReport)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, _ := json.Marshal(replyReport)
	_, _ = w.Write(data)
}
