package main

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func rdbConnect(addr string) *redis.Client {
	ctx := context.Background()
	log.Info().
		Str("redis", addr).
		Msg("connecting to redis...")
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
			log.Fatal().
				Err(err).
				Msg("Couldn't connect to redis")
			os.Exit(1)
		}
	}

	return rdb
}
