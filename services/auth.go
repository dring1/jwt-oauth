package services

import (
	"encoding/json"
	"net/http"

	"github.com/dring1/orm/lib/authentication"
	"github.com/dring1/orm/models"
)

func Login(requestUser *models.User) (int, []byte) {
	authBackend, err := authentication.JWTBackend()

	if authBackend.Authenticate(requestUser) {
		token, err := authBackend.GenerateToken(requestUser.UUID)
		if err != nil {
			return http.StatusInternalServerError, []byte("")
		} else {
			response, _ := json.Marshal(parameters.TokenAuthentication{token})
			return http.StatusOK, response
		}
	}

	return http.StatusUnauthorized, []byte("")
}

func RefreshToken(requestUser *models.User) []byte {
	jwtBackend := authentication.JWTBackend()
	token, err := jwtBackend.GenerateToken(requestUser.UUID)
}
