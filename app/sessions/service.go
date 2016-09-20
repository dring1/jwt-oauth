package sessions

import (
	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/jwt"
)

type Service interface {
	NewSession(string) (*Token, error)
	DeleteSession(string) error
}

type Session struct {
	Cache        *cache.Service
	TokenService jwt.Service
}

type Token struct {
	Value     string `json:"value"`
	UserEmail string `json:"email"`
}

func NewService(t jwt.Service, c *cache.Service) (Service, error) {
	return &Session{
		Cache:        c,
		TokenService: t,
	}, nil
}

// New - A new session is creating a token with correct TTL
// And scope
func (s *Session) NewSession(key string) (*Token, error) {
	t, err := s.TokenService.NewToken(key)
	if err != nil {
		return nil, err
	}
	return &Token{
		UserEmail: key,
		Value:     t,
	}, nil
}
func (s *Session) DeleteSession(key string) error {

	return nil
}
