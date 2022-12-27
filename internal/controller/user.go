package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/feeds"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/util"
	"net/http"
	"time"
)

// GetUserPostsRSSV1 - get posts by a user in rss format
func GetUserPostsRSSV1(c *gin.Context) {
	c.Header("Cache-Control", "max-age=30")
	username := c.Param("username")
	session, _ := util.GetSession(c)
	user, err := service.CreateUserService().GetUserByUsername(username)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	nameToShow := "@" + username
	if user.Name != "" {
		nameToShow = user.Name + " (" + nameToShow + ")"
	}
	posts, _ := service.CreatePostService().GetPostsForUser(
		session,
		username,
		constants.UserPostsDefaultPageSize,
	)
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
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetUserPostsV1 - get posts by a user
func GetUserPostsV1(c *gin.Context) {
	c.Header("Cache-Control", "max-age=30")
	username := c.Param("username")
	session, _ := util.GetSession(c)
	posts, _ := service.CreatePostService().GetPostsForUser(
		session,
		username,
		constants.UserPostsDefaultPageSize,
	)
	c.JSON(http.StatusOK, posts)
}

// GetSuggestedFollowsForUserV1 - Get suggested follows for user
func GetSuggestedFollowsForUserV1(c *gin.Context) {
	uuidParam, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	users := service.CreateUserService().GetSuggestedFollowsForUser(uuidParam)
	c.JSON(http.StatusOK, users)
}
