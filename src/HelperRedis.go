package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func redisConnect() {
	// Create Redis Client
	client := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_URL", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}
