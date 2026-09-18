package routes

import (
	"github.com/gin-gonic/gin"
	"live-polling-app/controllers"
	"live-polling-app/middleware"
)

func Setup(r *gin.Engine, auth *controllers.AuthController, poll *controllers.PollController, secret string) {
	api := r.Group("/api")

	api.POST("/auth/signup", auth.Signup)
	api.POST("/auth/login", auth.Login)

	api.POST("/polls", middleware.RequireAuth(secret), poll.CreatePoll)
	api.GET("/polls/my", middleware.RequireAuth(secret), poll.MyPolls)

	api.GET("/polls/:id", poll.GetPoll)
	api.GET("/polls/:id/results", poll.Results)
	api.POST("/polls/:id/vote", poll.Vote)
	api.GET("/polls/:id/live", poll.Live)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
