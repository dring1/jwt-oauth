package cache

import "gopkg.in/redis.v3"

type Service struct {
	*redis.Client
}

func New(redisEndpoint string) (*Service, error) {
	client := Service{redis.NewClient(&redis.Options{
		Addr:     redisEndpoint,
		Password: "", // no password set
		DB:       0,  // use default DB
	})}

	return &client, nil
}
