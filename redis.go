package main

import (
	"github.com/go-redis/redis/v8"
)

func rdbConnect() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.176.2:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
