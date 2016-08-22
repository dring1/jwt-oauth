package jwtbackend

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	privateKey, PublicKey []byte
	TokenDuration         int
	ExpireOffset          int
}

const (
	TokenDuration = 72
	ExpireOffset  = 3600
)

func NewJWTBackend(privateKey []byte, publicKey []byte, tokenDuration, expireOffset int) (*JWTService, error) {
	return &JWTService{
		privateKey:    privateKey,
		PublicKey:     publicKey,
		TokenDuration: TokenDuration,
		ExpireOffset:  expireOffset,
	}, nil
}

type CustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func (backend *JWTService) GenerateToken(userID string) (string, error) {
	exp := time.Now().Add(time.Hour * 1).Unix()
	iss := "jwt-oauth.com"
	sub := "jwt-oauth"
	claims := CustomClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: exp,
			Issuer:    iss,
			Subject:   sub,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Insert into cache here ?
func (backend *JWTService) Authenticate(interface{}) bool {
	return true
}

func (backend *JWTService) Logout(tokenString string, token *jwt.Token) error {
	return nil
}

func (backend *JWTService) TimeToExpire(timestamp interface{}) int64 {

	if ts, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(ts), 0)
		if remainder := tm.Sub(time.Now()); remainder > 0 {
			return int64(remainder.Seconds()) + ExpireOffset
		}
	}
	return ExpireOffset
}
