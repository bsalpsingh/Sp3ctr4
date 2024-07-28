package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUsers(c *gin.Context) {
	fmt.Println("URL : ", c.Request.URL)
	c.JSON(http.StatusOK, gin.H{
		"data": "welcome to users space",
	})
}
