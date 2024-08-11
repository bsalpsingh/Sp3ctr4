package middlewares

import (
	"net/http"
	"strings"

	"github.com/sp3ctr4/auth"
	"github.com/sp3ctr4/database"

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
				c.JSON(http.StatusInternalServerError, gin.H{
					"data": "Failed to authorize user",
				})
				c.Abort()
				return
			}

			var authorizedUser database.User
			email := claims["email"].(string)
			if query := database.DB.Find(&authorizedUser, "email = ? ", email); query.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"data": "user not found",
				})
				c.Abort()
				return
			}
			c.Set("user", authorizedUser)

			c.Next()
			return

		}

	}
}
