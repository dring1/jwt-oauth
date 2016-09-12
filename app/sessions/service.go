package sessions

import "github.com/dring1/jwt-oauth/cache"

type Service interface {
	New() error
	Delete() error
}

type S struct {
	Cache *cache.Service
}

func (s *S) New(key string) error {

	return nil
}
