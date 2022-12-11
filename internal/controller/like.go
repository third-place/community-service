package controller

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/third-place/community-service/internal/service"
	"net/http"
)

// CreateNewPostLikeV1 - like a post
func CreateNewPostLikeV1(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := uuid.MustParse(params["uuid"])
	session, err := service.CreateDefaultAuthService().GetSessionFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	postLike, err := service.CreateDefaultLikeService().CreateLikeForPost(uuidParam, uuid.MustParse(session.User.Uuid))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, _ := json.Marshal(postLike)
	_, _ = w.Write(data)
}

// DeleteLikeForPostV1 - delete like for post
func DeleteLikeForPostV1(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := uuid.MustParse(params["uuid"])
	session, err := service.CreateDefaultAuthService().GetSessionFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err = service.CreateDefaultLikeService().DeleteLikeForPost(uuidParam, uuid.MustParse(session.User.Uuid))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}
