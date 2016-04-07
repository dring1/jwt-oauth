package services

import (
	"net/http"

	"github.com/dring1/orm/models"
)

func Login(requestUser *models.User) (int, []byte) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	//
	// if authBackend.Authenticate(requestUser) {
	//     token, err := authBackend.GenerateToken(requestUser.UUID)
	//     if err != nil {
	//         return http.StatusInternalServerError, []byte("")
	//     } else {
	//         response, _ := json.Marshal(parameters.TokenAuthentication{token})
	//         return http.StatusOK, response
	//     }
	// }
	//
	return http.StatusUnauthorized, []byte("")
}
