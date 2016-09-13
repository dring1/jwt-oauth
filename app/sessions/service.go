package sessions

import "github.com/dring1/jwt-oauth/cache"

type Service interface {
	New(string) (*Token, error)
	Delete(string) error
}

type Session struct {
	Cache *cache.Service
}

type Token struct {
	Value     string `json:"value"`
	UserEmail string `json:"email"`
}

func NewService(c *cache.Service) Service {
	return &Session{
		c,
	}
}

// New - A new session is creating a token with correct TTL
// And scope
func (s *Session) New(key string) (*Token, error) {

	return nil, nil
}
func (s *Session) Delete(key string) error {

	return nil
}
