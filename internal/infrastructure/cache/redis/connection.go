package redis

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func ConnectionCache() *redis.Client {
	godotenv.Load()

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Panicf("redis connection error: %v", err)
	}

	log.Println("redis connected")
	return client
}