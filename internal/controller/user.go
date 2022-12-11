package controller

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/util"
	uuid2 "github.com/third-place/community-service/internal/uuid"
	"net/http"
	"time"
)

// GetUserPostsRSSV1 - get posts by a user in rss format
func GetUserPostsRSSV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=30")
	params := mux.Vars(r)
	username := params["username"]
	session, _ := util.GetSession(r.Header.Get("x-session-token"))
	var viewerUuid uuid.UUID
	if session != nil {
		viewerUuid = uuid.MustParse(session.User.Uuid)
	}
	user, err := service.CreateDefaultUserService().GetUserByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	nameToShow := "@" + username
	if user.Name != "" {
		nameToShow = user.Name + " (" + nameToShow + ")"
	}
	posts, _ := service.CreatePostService().GetPostsForUser(
		username, &viewerUuid, constants.UserPostsDefaultPageSize)
	feed := &feeds.Feed{
		Title:       "RSS for @" + username + " - Third place",
		Link:        &feeds.Link{Href: "https://thirdplaceapp.com/posts/" + username + "/rss"},
		Description: "Posts by @" + username + ", provided by Third place.",
		Author:      &feeds.Author{Name: nameToShow},
		Created:     time.Now(),
	}
	var feedItems []*feeds.Item
	for _, post := range posts {
		feedItems = append(feedItems, &feeds.Item{
			Id:          post.Uuid,
			Link:        &feeds.Link{Href: "https://thirdplaceapp.com/p/" + post.Uuid},
			Description: post.Text,
			Created:     post.CreatedAt,
		})
	}
	feed.Items = feedItems
	data, err := feed.ToRss()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(data))
}

// GetUserPostsV1 - get posts by a user
func GetUserPostsV1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "max-age=30")
	params := mux.Vars(r)
	username := params["username"]
	session, _ := util.GetSession(r.Header.Get("x-session-token"))
	var viewerUuid uuid.UUID
	if session != nil {
		viewerUuid = uuid.MustParse(session.User.Uuid)
	}
	posts, _ := service.CreatePostService().GetPostsForUser(
		username, &viewerUuid, constants.UserPostsDefaultPageSize)
	data, _ := json.Marshal(posts)
	_, _ = w.Write(data)
}

// GetSuggestedFollowsForUserV1 - Get suggested follows for user
func GetSuggestedFollowsForUserV1(w http.ResponseWriter, r *http.Request) {
	users := service.CreateDefaultUserService().
		GetSuggestedFollowsForUser(uuid2.GetUuidFromPathSecondPosition(r.URL.Path))
	data, _ := json.Marshal(users)
	_, _ = w.Write(data)
}
