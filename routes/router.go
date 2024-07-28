package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sp3ctr4/utils"
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

	authGrp := v1.Group("/auth")

	authGrp.GET("/login", utils.HandleGoogleLogin)
	authGrp.GET("/callback", utils.HandleGoogleCallback)

	usersGrp.GET("/", getUsers)

}
