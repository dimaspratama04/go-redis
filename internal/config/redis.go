package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		//TLSConfig: &tls.Config{},
		DB: 0})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal("failed to connect redis:", err)
	}

	log.Println("Connected to redis")

	return rdb
}
