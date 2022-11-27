package controller

import (
	model2 "github.com/third-place/community-service/internal/auth/model"
	"github.com/third-place/community-service/internal/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

// CreateNewPostLikeV1 - like a post
func CreateNewPostLikeV1(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := uuid.MustParse(params["uuid"])

	service.CreateDefaultAuthService().DoWithValidSession(w, r, func(session *model2.Session) (interface{}, error) {
		return service.CreateDefaultLikeService().CreateLikeForPost(uuidParam, uuid.MustParse(session.User.Uuid))
	})
}

// DeleteLikeForPostV1 - delete like for post
func DeleteLikeForPostV1(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uuidParam := uuid.MustParse(params["uuid"])

	service.CreateDefaultAuthService().DoWithValidSession(w, r, func(session *model2.Session) (interface{}, error) {
		_ = service.CreateDefaultLikeService().DeleteLikeForPost(uuidParam, uuid.MustParse(session.User.Uuid))
		return nil, nil
	})

}
