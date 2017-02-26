package token

import (
	"fmt"
	"time"

	_jwt "github.com/dgrijalva/jwt-go"
	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/lib/errors"
	"github.com/dring1/jwt-oauth/models"
)

type TokenService struct {
	cache                 *cache.Service
	privateKey, PublicKey []byte
	TokenTTL              time.Duration
	ExpireOffset          int
	TokenISS              string
	TokenSub              string
}

type TimeStamp int64

type Service interface {
	NewToken(string) (*Token, error)
	RefreshToken(*Token) (*Token, error)
	TimeToExpire(TimeStamp) TimeStamp
	Validate(string) (*Token, bool, error)
	Revoke(*Token) error
	IsRevoked(*Token) (bool, error)
}

type CustomClaims struct {
	Email string `json:"email"`
	_jwt.StandardClaims
}

type Token models.Token

func NewService(privKey, publicKey []byte, tokenTTL int, expireOffset int, tokISS, tokSub string, cache *cache.Service) (Service, error) {
	return &TokenService{
		cache:        cache,
		privateKey:   privKey,
		PublicKey:    publicKey,
		TokenTTL:     time.Duration(tokenTTL) * time.Second,
		ExpireOffset: expireOffset,
		TokenISS:     tokISS,
		TokenSub:     tokSub,
	}, nil

}

func (t *TokenService) NewToken(userID string) (*Token, error) {
	exp := time.Now().Add(t.TokenTTL).Unix()
	iss := t.TokenISS
	sub := t.TokenSub
	claims := CustomClaims{
		userID,
		_jwt.StandardClaims{
			ExpiresAt: exp,
			Issuer:    iss,
			Subject:   sub,
		},
	}
	token := _jwt.NewWithClaims(_jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.privateKey)
	if err != nil {
		return &Token{}, err
	}
	return &Token{T: *token, TokenString: tokenString}, nil
}

func (ts *TokenService) Validate(tokenString string) (*Token, bool, error) {
	token, err := ts.parseToken(tokenString)
	if err != nil {
		return nil, false, err
	}
	revoked, err := ts.IsRevoked(token)
	if revoked || err != nil {
		return nil, false, err
	}
	return token, token.T.Valid, nil
}
func (t *TokenService) TimeToExpire(timestamp TimeStamp) TimeStamp {
	tm := time.Unix(int64(timestamp), 0)
	if remainder := tm.Sub(time.Now()); remainder > 0 {
		return TimeStamp(int(remainder.Seconds()) + t.ExpireOffset)
	}
	return 0
}

func (t *TokenService) parseToken(tokenString string) (*Token, error) {
	tok, err := _jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *_jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*_jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(errors.InvalidToken)
		}
		return []byte(t.privateKey), nil
	})
	return &Token{T: *tok}, err
	//return (*Token)(tok), err
}

func (t *TokenService) Revoke(token *Token) error {
	// the duration of the token
	return t.cache.Set(token.T.Raw, 0, t.TokenTTL).Err()
}

func (t *TokenService) IsRevoked(token *Token) (bool, error) {
	return t.cache.Exists(token.T.Raw).Result()
}

func (t *TokenService) RefreshToken(token *Token) (*Token, error) {
	// revoke token
	err := t.Revoke(token)
	if err != nil {
		return &Token{}, err
	}
	claims, ok := token.T.Claims.(*CustomClaims)
	if !ok {
		return &Token{}, fmt.Errorf(errors.InvalidToken)
	}
	return t.NewToken(claims.Email)
}
