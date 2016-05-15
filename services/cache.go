package services

import "gopkg.in/redis.v3"

var Client *redis.Client

func NewCacheClient() *redis.Client {
	Client = redis.NewClient(&redis.Options{})
	return Client
}
