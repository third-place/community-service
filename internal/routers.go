/*
 * Otto user service
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package internal

import (
	"github.com/gorilla/mux"
	"github.com/third-place/community-service/internal/controller"
	"net/http"
	"strings"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	{
		"CreateAReplyV1",
		strings.ToUpper("Post"),
		"/reply",
		controller.CreateAReplyV1,
	},
	{
		"CreateAReplyReportV1",
		strings.ToUpper("Post"),
		"/report/reply",
		controller.CreateReplyReportV1,
	},
	{
		"CreateAPostReport",
		strings.ToUpper("Post"),
		"/report",
		controller.CreatePostReportV1,
	},
	{
		"CreateAShare",
		strings.ToUpper("Post"),
		"/share",
		controller.CreateShareV1,
	},
	{
		"CreateFollowV1",
		strings.ToUpper("Post"),
		"/follow",
		controller.CreateFollowV1,
	},
	{
		"GetSharesV1",
		strings.ToUpper("Get"),
		"/share",
		controller.GetSharesV1,
	},
	{
		"CreateNewPostLikeV1",
		strings.ToUpper("Post"),
		"/post/{uuid}/like",
		controller.CreateNewPostLikeV1,
	},
	{
		"CreateNewPostV1",
		strings.ToUpper("Post"),
		"/post",
		controller.CreateNewPostV1,
	},
	{
		"DeleteFollowV1",
		strings.ToUpper("Delete"),
		"/follow/{uuid}",
		controller.DeleteFollowV1,
	},
	{
		"DeleteLikeForPostV1",
		strings.ToUpper("Delete"),
		"/post/{uuid}/like",
		controller.DeleteLikeForPostV1,
	},
	{
		"DeletePostV1",
		strings.ToUpper("Delete"),
		"/post/{uuid}",
		controller.DeletePostV1,
	},
	{
		"GetDraftPostsV1",
		strings.ToUpper("Get"),
		"/post/draft",
		controller.GetDraftPostsV1,
	},
	{
		"GetLikedPostsV1",
		strings.ToUpper("Get"),
		"/likes/{username}",
		controller.GetLikedPostsV1,
	},
	{
		"GetPostsFirehoseV1",
		strings.ToUpper("Get"),
		"/post",
		controller.GetPostsFirehoseV1,
	},
	{
		"GetPostV1",
		strings.ToUpper("Get"),
		"/post/{uuid}",
		controller.GetPostV1,
	},
	{
		"GetPostRepliesV1",
		strings.ToUpper("Get"),
		"/post/{uuid}/replies",
		controller.GetPostRepliesV1,
	},
	{
		"GetShareV1",
		strings.ToUpper("Get"),
		"/share/{uuid}",
		controller.GetShareV1,
	},
	{
		"GetSuggestedFollowsForUserV1",
		strings.ToUpper("Get"),
		"/suggested-follows/{user}",
		controller.GetSuggestedFollowsForUserV1,
	},
	{
		"GetUserFollowersV1",
		strings.ToUpper("Get"),
		"/followers/{username}",
		controller.GetUserFollowersV1,
	},
	{
		"GetUserFollowsV1",
		strings.ToUpper("Get"),
		"/follows/{username}",
		controller.GetUserFollowsV1,
	},
	{
		"GetPostsForUserFollowsV1",
		strings.ToUpper("Get"),
		"/follow-posts/{username}",
		controller.GetPostsForUserFollowsV1,
	},
	{
		"GetUserPostsRSSV1",
		strings.ToUpper("Get"),
		"/posts/{username}/rss",
		controller.GetUserPostsRSSV1,
	},
	{
		"GetUserPostsV1",
		strings.ToUpper("Get"),
		"/posts/{username}",
		controller.GetUserPostsV1,
	},
	{
		"GetNewPostsV1",
		strings.ToUpper("Get"),
		"/new-posts/{username}",
		controller.GetNewPostsV1,
	},
	{
		"UpdatePostV1",
		strings.ToUpper("Put"),
		"/post",
		controller.UpdatePostV1,
	},
}
