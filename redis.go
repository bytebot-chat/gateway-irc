package main

import (
	"github.com/go-redis/redis/v8"
)

func rdbConnect(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
