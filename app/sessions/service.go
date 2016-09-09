package sessions

import "github.com/dring1/jwt-oauth/cache"

type Session interface {
	Create() error
	Delete() error
}

type Service struct {
	Cache *cache.Service
}

func (s *Service) Create() error {
	// TODO: Push TTL
	return nil
}
