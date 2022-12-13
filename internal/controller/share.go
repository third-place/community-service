package controller

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/util"
	"net/http"
)

// GetShareV1 - get a share
func GetShareV1(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := uuid.MustParse(params["uuid"])
	share, err := service.CreateShareService().GetShare(uuidParam)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, _ := json.Marshal(share)
	_, _ = w.Write(data)
}

// GetSharesV1 - get shares
func GetSharesV1(w http.ResponseWriter, r *http.Request) {
	session, _ := util.GetSession(r.Header.Get("x-session-token"))
	var viewerUsername string
	if session != nil {
		viewerUser, _ := service.CreateUserService().GetUser(uuid.MustParse(session.User.Uuid))
		viewerUsername = viewerUser.Username
	}
	limit := constants.UserPostsDefaultPageSize
	share, err := service.CreateShareService().GetShares(&viewerUsername, limit)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, _ := json.Marshal(share)
	_, _ = w.Write(data)
}

// CreateShareV1 - create a reshare
func CreateShareV1(w http.ResponseWriter, r *http.Request) {
	newShareParam := model.DecodeRequestToNewShare(r)
	share, err := service.CreateShareService().CreateShare(newShareParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, _ := json.Marshal(share)
	_, _ = w.Write(data)
}
