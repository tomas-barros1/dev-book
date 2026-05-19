package routes

import (
	"dev-book-api/handlers"
	"dev-book-api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(h *handlers.Handler, g *gin.Engine, jwtSecret string) {
	g.POST("/login", h.LoginUser)

	users := g.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("", h.SearchUsers)
		users.GET("/:userId", h.GetUser)

		authenticated := users.Group("")
		authenticated.Use(middleware.AuthMiddleware(jwtSecret))
		{
			authenticated.PUT("/:userId", h.UpdateUser)
			authenticated.DELETE("/:userId", h.DeleteUser)
			authenticated.POST("/:userId/follow", h.FollowUser)
			authenticated.POST("/:userId/unfollow", h.UnfollowUser)
			authenticated.GET("/:userId/followers", h.GetFollowers)
			authenticated.GET("/:userId/following", h.GetFollowing)
			authenticated.POST("/:userId/update-password", h.UpdatePassword)
			authenticated.GET("/:userId/posts", h.GetUserPosts)
		}
	}

	posts := g.Group("/posts")
	posts.Use(middleware.AuthMiddleware(jwtSecret))
	{
		posts.POST("", h.CreatePost)
		posts.GET("", h.GetFeed)
		posts.GET("/:postId", h.GetPost)
		posts.PUT("/:postId", h.UpdatePost)
		posts.DELETE("/:postId", h.DeletePost)
		posts.POST("/:postId/like", h.LikePost)
		posts.POST("/:postId/unlike", h.UnlikePost)
	}
}
