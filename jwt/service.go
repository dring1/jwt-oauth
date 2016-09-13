package jwt

import (
	"time"

	_jwt "github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	privateKey, PublicKey []byte
	TokenDuration         int
	ExpireOffset          int
	TokenISS              string
	TokenSub              string
}

type TimeStamp int64

type Service interface {
	NewToken()
	TimeToExpire(TimeStamp)
}

type CustomClaims struct {
	Email string `json:"email"`
	_jwt.StandardClaims
}

func New() (*Service, error) {
	return nil, nil

}

func (backend *JWTService) GenerateToken(userID string) (string, error) {
	exp := time.Now().Add(time.Duration(backend.ExpireOffset)).Unix()
	iss := backend.TokenISS
	sub := backend.TokenSub
	claims := CustomClaims{
		userID,
		_jwt.StandardClaims{
			ExpiresAt: exp,
			Issuer:    iss,
			Subject:   sub,
		},
	}
	token := _jwt.NewWithClaims(_jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Insert into cache here ?
// func (backend *JWTService) Authenticate(interface{}) bool {
// 	return true
// }

// func (backend *JWTService) Logout(tokenString string, token *jwt.Token) error {
// 	return nil
// }

func (backend *JWTService) TimeToExpire(timestamp TimeStamp) TimeStamp {

	tm := time.Unix(int64(timestamp), 0)
	if remainder := tm.Sub(time.Now()); remainder > 0 {
		return TimeStamp(int(remainder.Seconds()) + backend.ExpireOffset)
	}
	return 0
}
