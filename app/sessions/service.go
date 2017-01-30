package sessions

import (
	"github.com/dring1/jwt-oauth/cache"
	"github.com/dring1/jwt-oauth/token"
)

type Service interface {
	NewSession(string) (*Token, error)
	DeleteSession(string) error
}

type Session struct {
	Cache        *cache.Service
	TokenService token.Service
}

type Token struct {
	Value string `json:"value"`
	Email string `json:"email"`
}

func NewService(t token.Service, c *cache.Service) (Service, error) {
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
		Email: key,
		Value: t,
	}, nil
}
func (s *Session) DeleteSession(key string) error {

	return nil
}
