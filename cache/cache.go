package cache

import "gopkg.in/redis.v3"

type CacheService struct {
	*redis.Client
}

func NewCacheService(redisEndpoint string) (*CacheService, error) {
	client := CacheService{redis.NewClient(&redis.Options{
		Addr:     redisEndpoint,
		Password: "", // no password set
		DB:       0,  // use default DB
	})}

	return &client, nil
}
