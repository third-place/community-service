package controller

import (
	"encoding/json"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/util"
	"net/http"
)

// CreatePostReportV1 - report a post
func CreatePostReportV1(w http.ResponseWriter, r *http.Request) {
	_, err := util.GetSession(r.Header.Get("x-session-token"))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	newReport := model.DecodeRequestToNewPostReport(r)
	report, err := service.CreateReportService().CreatePostReport(newReport)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, _ := json.Marshal(report)
	_, _ = w.Write(data)
}

// CreateReplyReportV1 - report a reply
func CreateReplyReportV1(w http.ResponseWriter, r *http.Request) {
	_, err := util.GetSession(r.Header.Get("x-session-token"))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	newReport := model.DecodeRequestToNewPostReport(r)
	replyReport, err := service.CreateReportService().CreateReplyReport(newReport)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, _ := json.Marshal(replyReport)
	_, _ = w.Write(data)
}
