package controller

import (
	"encoding/json"
	"github.com/danielmunro/otto-community-service/internal/model"
	"github.com/danielmunro/otto-community-service/internal/service"
	iUuid "github.com/danielmunro/otto-community-service/internal/uuid"
	"net/http"
)

// CreateAReplyV1 - create a reply
func CreateAReplyV1(w http.ResponseWriter, r *http.Request) {
	session := service.CreateDefaultAuthService().GetSessionFromRequest(r)
	if session == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	newReplyModel := model.DecodeRequestToNewReply(r)
	reply, err := service.CreateReplyService().CreateReply(session, newReplyModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	data, _ := json.Marshal(reply)
	_, _ = w.Write(data)
}

func GetPostRepliesV1(w http.ResponseWriter, r *http.Request) {
	postUuid := iUuid.GetUuidFromPathSecondPosition(r.URL.Path)
	replies, err := service.CreateReplyService().GetRepliesForPost(postUuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, _ := json.Marshal(replies)
	_, _ = w.Write(data)
}
