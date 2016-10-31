package token

import (
	"fmt"
	"time"

	_jwt "github.com/dgrijalva/jwt-go"
	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/lib/errors"
)

type TokenService struct {
	cache                 cache.Service
	privateKey, PublicKey []byte
	TokenTTL              int
	ExpireOffset          int
	TokenISS              string
	TokenSub              string
}

type TimeStamp int64

type Service interface {
	NewToken(string) (string, error)
	TimeToExpire(TimeStamp) TimeStamp
	Validate(string) (bool, error)
	InvalidateToken(string) error
}

type CustomClaims struct {
	Email string `json:"email"`
	_jwt.StandardClaims
}

func NewService(privKey, publicKey []byte, tokenTTL int, expireOffset int, tokISS, tokSub string) (Service, error) {
	return &TokenService{
		privateKey:   privKey,
		PublicKey:    publicKey,
		TokenTTL:     tokenTTL,
		ExpireOffset: expireOffset,
		TokenISS:     tokISS,
		TokenSub:     tokSub,
	}, nil

}

func (backend *TokenService) NewToken(userID string) (string, error) {
	exp := time.Now().Add(5 * time.Minute).Unix()
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
// func (backend *TokenService) Authenticate(interface{}) bool {
// 	return true
// }

// func (backend *TokenService) Logout(tokenString string, token *jwt.Token) error {
// 	return nil
// }

func (ts *TokenService) Validate(tokenString string) (bool, error) {
	token, err := _jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *_jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*_jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(errors.InvalidToken)
		}
		// fmt.Println(len(strings.Split(tokenString, ".")))
		return []byte(ts.privateKey), nil
	})
	// if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
	//     fmt.Printf("hi %+v %v\n", claims, claims.StandardClaims.ExpiresAt)
	// } else {
	//     fmt.Println(err)
	// }
	// if err != nil || !token.Valid {
	// 	return false, err
	// }
	fmt.Println(token.Method.(*_jwt.SigningMethodHMAC), token.Header)
	return token.Valid, err
}
func (backend *TokenService) TimeToExpire(timestamp TimeStamp) TimeStamp {

	tm := time.Unix(int64(timestamp), 0)
	if remainder := tm.Sub(time.Now()); remainder > 0 {
		return TimeStamp(int(remainder.Seconds()) + backend.ExpireOffset)
	}
	return 0
}

func (t *TokenService) parseToken(tokenString string) (*_jwt.Token, error) {

	return _jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *_jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*_jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(errors.InvalidToken)
		}
		// fmt.Println(len(strings.Split(tokenString, ".")))
		return []byte(t.privateKey), nil
	})
}

func (t *TokenService) InvalidateToken(token string) error {
	return nil
}
