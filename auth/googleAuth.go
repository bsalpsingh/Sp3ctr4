// oauth/oauth.go

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofor-little/env"
	"github.com/sp3ctr4/database"
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
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
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

	email, emailOk := userInfo["email"].(string)
	givenName, givenNameOk := userInfo["given_name"].(string)
	familyName, familyNameOk := userInfo["family_name"].(string)
	if !emailOk || !givenNameOk || !familyNameOk {
		log.Println("missing user information")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	// check if the user is presaent in the db, if no add user to db and sign token else query and sign token

	var newUser database.User

	if userQeuery := database.DB.Where(" email = ? ", userInfo["email"]).First(&newUser); userQeuery.Error != nil {

		newRegisteredUser := database.User{
			Name:  fmt.Sprintf("%v %v", givenName, familyName),
			Email: email,
		}
		if userRegisterQuery := database.DB.Create(&newRegisteredUser); userRegisterQuery.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"data": "failed to signup user",
			})
			c.Abort()
			return
		}
		newUser = newRegisteredUser
	}

	if signedToken, err := signToken(newUser); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"data": err.Error(),
		})
		return
	} else {

		c.IndentedJSON(http.StatusOK, gin.H{
			"data": gin.H{
				"token": signedToken,
				"user":  newUser,
			},
		})
	}

}
