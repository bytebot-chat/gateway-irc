package main

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	log "gopkg.in/inconshreveable/log15.v2"
)

func rdbConnect(addr string) *redis.Client {
	ctx := context.Background()
	log.Info("connecting to redis...", "redis", addr)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		time.Sleep(3 * time.Second)
		err := rdb.Ping(ctx).Err()
		if err != nil {
			log.Crit("FATAL unable to connect to redis", "error", err)
			os.Exit(1)
		}
	}

	return rdb
}
