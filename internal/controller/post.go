package controller

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/util"
	iUuid "github.com/third-place/community-service/internal/uuid"
	"net/http"
)

// CreateNewPostV1 - create a new post
func CreateNewPostV1(w http.ResponseWriter, r *http.Request) {
	session, err := util.GetSession(r.Header.Get("x-session-token"))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	newPostModel, _ := model.DecodeRequestToNewPost(r)
	post, err := service.CreatePostService().CreatePost(session, newPostModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data, _ := json.Marshal(post)
	_, _ = w.Write(data)
}

// UpdatePostV1 - update a post
func UpdatePostV1(w http.ResponseWriter, r *http.Request) {
	postModel, _ := model.DecodeRequestToPost(r)
	session, err := util.GetSession(r.Header.Get("x-session-token"))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err = service.CreatePostService().UpdatePost(session, postModel)
	if err != nil {
		w.WriteHeader(400)
	}
}

// GetPostV1 - get a post
func GetPostV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=60")
	session, _ := util.GetSession(r.Header.Get("x-session-token"))
	post, err := service.CreatePostService().GetPost(
		session,
		iUuid.GetUuidFromPathSecondPosition(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, _ := json.Marshal(post)
	_, _ = w.Write(data)
}

// GetPostsForUserFollowsV1 - get a user's friend's posts
func GetPostsForUserFollowsV1(w http.ResponseWriter, r *http.Request) {
	limit := constants.UserPostsDefaultPageSize
	params := mux.Vars(r)
	username := params["username"]
	session, _ := util.GetSession(r.Header.Get("x-session-token"))
	var viewerUuid uuid.UUID
	if session != nil {
		viewerUuid = uuid.MustParse(session.User.Uuid)
	}
	posts, err := service.CreatePostService().GetPostsForUserFollows(username, viewerUuid, limit)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, _ := json.Marshal(posts)
	_, _ = w.Write(data)
}

// GetDraftPostsV1 - get draft posts
func GetDraftPostsV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=30")
	session, err := util.GetSession(r.Header.Get("x-session-token"))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	posts := service.CreatePostService().GetDraftPosts(
		uuid.MustParse(session.User.Uuid),
		constants.UserPostsDefaultPageSize,
	)
	data, _ := json.Marshal(posts)
	_, _ = w.Write(data)
}

// GetPostsFirehoseV1 - get posts
func GetPostsFirehoseV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=30")
	session, err := util.GetSession(r.Header.Get("x-session-token"))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var viewerUsername string
	if session != nil {
		viewerUser, _ := service.CreateUserService().GetUser(uuid.MustParse(session.User.Uuid))
		viewerUsername = viewerUser.Username
	}
	limit := constants.UserPostsDefaultPageSize
	var posts []*model.Post
	posts, _ = service.CreatePostService().GetPostsFirehose(&viewerUsername, limit)
	data, _ := json.Marshal(posts)
	_, _ = w.Write(data)
}

// GetLikedPostsV1 - get liked posts
func GetLikedPostsV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=30")
	params := mux.Vars(r)
	username := params["username"]
	limit := constants.UserPostsDefaultPageSize
	var posts []*model.Post
	posts, _ = service.CreatePostService().GetLikedPosts(username, limit)
	data, _ := json.Marshal(posts)
	_, _ = w.Write(data)
}

// DeletePostV1 - delete a post
func DeletePostV1(w http.ResponseWriter, r *http.Request) {
	session, err := util.GetSession(r.Header.Get("x-session-token"))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	params := mux.Vars(r)
	postUuid, err := uuid.Parse(params["uuid"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = service.CreatePostService().DeletePost(session, postUuid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}
