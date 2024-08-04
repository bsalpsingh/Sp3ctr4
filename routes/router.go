package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sp3ctr4/auth"
	"github.com/sp3ctr4/middlewares"
)

func Init(e *gin.Engine) {
	v1 := e.Group("/api/v1")

	// home route

	v1.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": "welcome to sp3ctr4",
		})
	})

	// user route
	usersGrp := v1.Group("/users")

	usersGrp.Use(middlewares.IsLoggedIn())
	usersGrp.GET("/", getUsers)

	authGrp := v1.Group("/auth")

	authGrp.GET("/login", auth.HandleGoogleLogin)
	authGrp.GET("/callback", auth.HandleGoogleCallback)

}
