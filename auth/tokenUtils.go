package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sp3ctr4/database"
)

var (
	key []byte
	t   *jwt.Token
)

func signToken(user database.User) (string, error) {

	key = []byte("my_super_secret_key")
	t = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"name": user.Name,

			"id": user.ID,
		})
	return t.SignedString(key)
}

func ParseToken(tokenString string) (map[string]interface{}, error) {
	key = []byte("my_super_secret_key")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

		}
		return key, nil
	})

	if err != nil {

		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("failed to parse token")
	}
}
