package services

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dring1/orm/lib/authentication"
	"github.com/dring1/orm/models"
)

func Login(requestUser *models.User) ([]byte, error) {
	authBackend := authentication.JWTBackendInstance

	if ok := authBackend.Authenticate(requestUser); !ok {
		// return unauthorized
		return []byte(""), nil
	}
	token, err := authBackend.GenerateToken(requestUser.Email)
	if err != nil {
		return []byte(""), err
	}
	// Insert token into cache
	response, _ := json.Marshal(authentication.AuthToken{T: token})
	return response, nil

}

func RefreshToken(requestUser *models.User) ([]byte, error) {
	jwtBackend := authentication.JWTBackendInstance
	token, err := jwtBackend.GenerateToken(requestUser.Email)
	if err != nil {
		return nil, err
	}
	response, _ := json.Marshal(authentication.AuthToken{T: token})
	return response, nil
}

func Logout(req *http.Request) error {
	authBackend := authentication.JWTBackendInstance
	tokenRequest, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err != nil {
		return err
	}
	tokenString := req.Header.Get("Authorization")
	return authBackend.Logout(tokenString, tokenRequest)
}
