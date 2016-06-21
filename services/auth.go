package services

import (
	"encoding/json"
	"net/http"
)

func Login(key string) ([]byte, error) {
	authBackend := JWTBackend()

	if ok := JWTBackend().Authenticate(key); !ok {
		// return unauthorized
		return []byte(""), nil
	}
	token, err := authBackend.GenerateToken(key)
	if err != nil {
		return []byte(""), err
	}
	// Insert token into cache
	response, _ := json.Marshal(AuthToken{T: token})
	return response, nil

}

func RefreshToken(key string) ([]byte, error) {
	// jwtBackend := services.JWTBackend
	token, err := JWTBackend().GenerateToken(key)
	if err != nil {
		return nil, err
	}
	response, _ := json.Marshal(AuthToken{T: token})
	return response, nil
}

func Logout(req *http.Request) error {
	// authBackend := JWTBackend()
	// tokenRequest, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
	// 	return authBackend.PublicKey, nil
	// })
	// if err != nil {
	// 	return err
	// }
	// tokenString := req.Header.Get("Authorization")
	// return authBackend.Logout(tokenString, tokenRequest)
	return nil
}
