package weblib

import (
	"gateway/weblib/handlers"
	"gateway/weblib/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(service map[string]interface{}) *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.InitMiddleware(service))

	v1 := ginRouter.Group("/douyin")
	{
		//user
		user := v1.Group("/user")
		{
			user.POST("/register/", handlers.Register)
			user.POST("/login/", handlers.Login)
			user.GET("/", handlers.UserInfo)
		}

		//publish
		publish := v1.Group("/publish")
		{
			publish.POST("/action/", handlers.Publish)
			publish.GET("/list/", handlers.PublishList)
		}

		//feed
		feed := v1.Group("/feed")
		{
			feed.GET("/", handlers.Feed)
		}

		favorite := v1.Group("/favorite")
		{
			favorite.POST("/action/", handlers.FavoriteAction)
			favorite.GET("/list/", handlers.FavoriteList)
		}

		comment := v1.Group("/comment")
		{
			comment.POST("/action/", handlers.CommentAction)
			comment.GET("/list/", handlers.CommentList)
		}

		relation := v1.Group("/relation")
		{
			relation.POST("/action/", handlers.RelationAction)
			relation.GET("/follow/list/", handlers.FollowList)
			relation.GET("/follower/list/", handlers.FollowerList)
			relation.GET("/friend/list/", handlers.FriendList)
		}

		message := v1.Group("/message")
		{
			message.POST("/action/", handlers.MessageAction)
			message.GET("/chat/", handlers.MessageChat)
		}
	}
	return ginRouter
}
