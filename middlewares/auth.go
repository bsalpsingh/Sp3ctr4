package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sp3ctr4/auth"

	"github.com/gin-gonic/gin"
)

func IsLoggedIn() gin.HandlerFunc {

	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"data": "user not logged in",
			})
			c.Abort()
		}
		if tokenSlice := strings.Split((authorization), " "); len(tokenSlice[1]) == 0 || strings.ToLower(tokenSlice[0]) != "bearer" {
			c.JSON(http.StatusForbidden, gin.H{
				"data": "user not logged in",
			})
			c.Abort()
			return

		} else {
			claims, err := auth.ParseToken(tokenSlice[1])

			if err != nil {
				c.Set("user", claims)
				fmt.Println("user : ", claims)
				c.Next()
				return
			}

		}

	}
}
