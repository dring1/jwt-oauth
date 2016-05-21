package services

import (
	"log"
	"sync"

	"gopkg.in/redis.v3"
)

type cache struct {
	*redis.Client
}

var (
	once   sync.Once
	client *cache
)

func Cache() *cache {
	once.Do(func() {
		log.Println("Creating Cache Client...")
		client = &cache{redis.NewClient(&redis.Options{})}
	})
	return client
}
