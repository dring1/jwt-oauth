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
	onceCache sync.Once
	client    cache
)

func Cache() *cache {
	onceCache.Do(func() {
		log.Println("Creating Cache Client...")
		client = cache{redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})}
	})
	return &client
}
