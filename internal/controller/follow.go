package controller

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/third-place/community-service/internal/mapper"
	"github.com/third-place/community-service/internal/service"
	iUuid "github.com/third-place/community-service/internal/uuid"
	"log"
	"net/http"
)

// CreateFollowV1 - create a follow
func CreateFollowV1(w http.ResponseWriter, r *http.Request) {
	newFollowModel, err := mapper.DecodeRequestToNewFollow(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	session := service.CreateDefaultAuthService().GetSessionFromRequest(r)
	if session == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	follow, err := service.CreateFollowService().CreateFollow(uuid.MustParse(session.User.Uuid), newFollowModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(follow)
	_, _ = w.Write(data)
}

// GetUserFollowersByUsernameV1 - get user followers
func GetUserFollowersByUsernameV1(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	usernameParam := params["username"]

	follows, err := service.CreateFollowService().GetUserFollowersByUsername(usernameParam)
	if err != nil {
		log.Print("error received from get user follows :: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, _ := json.Marshal(follows)
	_, _ = w.Write(data)
}

// GetUserFollowersV1 - get user followers
func GetUserFollowersV1(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	follows, err := service.CreateFollowService().GetUserFollowers(username)
	if err != nil {
		log.Print("error received from get user follows :: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, _ := json.Marshal(follows)
	_, _ = w.Write(data)
}

// GetUserFollowsV1 - get user follows
func GetUserFollowsV1(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	follows, err := service.CreateFollowService().GetUserFollows(username)
	if err != nil {
		log.Print("error received from get user follows :: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, _ := json.Marshal(follows)
	_, _ = w.Write(data)
}

// DeleteFollowV1 - delete a follow
func DeleteFollowV1(w http.ResponseWriter, r *http.Request) {
	session := service.CreateDefaultAuthService().GetSessionFromRequest(r)
	if session == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	followUuid := iUuid.GetUuidFromPathSecondPosition(r.URL.Path)
	err := service.CreateFollowService().DeleteFollow(followUuid, uuid.MustParse(session.User.Uuid))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}
