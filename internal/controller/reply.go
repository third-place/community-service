package controller

import (
	"encoding/json"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	iUuid "github.com/third-place/community-service/internal/uuid"
	"net/http"
)

// CreateAReplyV1 - create a reply
func CreateAReplyV1(w http.ResponseWriter, r *http.Request) {
	session, err := service.CreateDefaultAuthService().GetSessionFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
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
