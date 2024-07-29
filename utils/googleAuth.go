// oauth/oauth.go

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofor-little/env"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthStateString = "random" // Can be any random string for validation
// HandleGoogleLogin redirects to the Google OAuth2 consent page

func getGoogleAuthConfig() *oauth2.Config {
	if err := env.Load("./.env"); err != nil {
		panic(err)
	}
	ClientID := env.Get("googleClientId", " ")
	ClientSecret := env.Get("googleClientSecret", " ")
	PORT := env.Get("PORT", "8080")

	var (
		googleOauthConfig = &oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			RedirectURL:  fmt.Sprintf("http://localhost:%v/api/v1/auth/callback", PORT),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		}
	)
	return googleOauthConfig
}
func HandleGoogleLogin(c *gin.Context) {

	googleOauthConfig := getGoogleAuthConfig()

	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleGoogleCallback handles the callback from Google OAuth2
func HandleGoogleCallback(c *gin.Context) {
	googleOauthConfig := getGoogleAuthConfig()

	if c.Query("state") != oauthStateString {
		log.Println("invalid oauth state")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("could not get token: %v\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Printf("could not create request: %v\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		log.Printf("could not decode response: %v\n", err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	c.JSON(http.StatusOK, gin.H{"userInfo": userInfo})
}
