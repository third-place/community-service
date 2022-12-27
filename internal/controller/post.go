package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/third-place/community-service/internal/constants"
	"github.com/third-place/community-service/internal/model"
	"github.com/third-place/community-service/internal/service"
	"github.com/third-place/community-service/internal/util"
	"net/http"
)

// CreateNewPostV1 - create a new post
func CreateNewPostV1(c *gin.Context) {
	session, err := util.GetSession(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	newPostModel, _ := model.DecodeRequestToNewPost(c.Request)
	post, err := service.CreatePostService().CreatePost(session, newPostModel)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusCreated, post)
}

// UpdatePostV1 - update a post
func UpdatePostV1(c *gin.Context) {
	postModel, _ := model.DecodeRequestToPost(c.Request)
	session, err := util.GetSession(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	err = service.CreatePostService().UpdatePost(session, postModel)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
}

// GetPostV1 - get a post
func GetPostV1(c *gin.Context) {
	c.Header("Cache-Control", "max-age=60")
	session, _ := util.GetSession(c)
	postUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	post, err := service.CreatePostService().GetPost(
		session,
		postUuid,
	)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, post)
}

// GetPostsForUserFollowsV1 - get a user's friend's posts
func GetPostsForUserFollowsV1(c *gin.Context) {
	limit := constants.UserPostsDefaultPageSize
	username := c.Param("username")
	session, _ := util.GetSession(c)
	posts, err := service.CreatePostService().GetPostsForUserFollows(session, username, limit)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, posts)
}

// GetDraftPostsV1 - get draft posts
func GetDraftPostsV1(c *gin.Context) {
	c.Header("Cache-Control", "max-age=30")
	session, err := util.GetSession(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	posts := service.CreatePostService().GetDraftPosts(session, constants.UserPostsDefaultPageSize)
	c.JSON(http.StatusOK, posts)
}

// GetPostsFirehoseV1 - get posts
func GetPostsFirehoseV1(c *gin.Context) {
	c.Header("Cache-Control", "max-age=30")
	session, err := util.GetSession(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	limit := constants.UserPostsDefaultPageSize
	var posts []*model.Post
	posts, _ = service.CreatePostService().GetPostsFirehose(session, limit)
	c.JSON(http.StatusOK, posts)
}

// GetLikedPostsV1 - get liked posts
func GetLikedPostsV1(c *gin.Context) {
	c.Header("Cache-Control", "max-age=30")
	session, _ := util.GetSession(c)
	username := c.Param("username")
	limit := constants.UserPostsDefaultPageSize
	var posts []*model.Post
	posts, _ = service.CreatePostService().GetLikedPosts(session, username, limit)
	c.JSON(http.StatusOK, posts)
}

// DeletePostV1 - delete a post
func DeletePostV1(c *gin.Context) {
	session, err := util.GetSession(c)
	if err != nil {
		c.Status(http.StatusForbidden)
		return
	}
	postUuid, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	err = service.CreatePostService().DeletePost(session, postUuid)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
}
