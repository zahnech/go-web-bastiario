package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisDatabase *redis.Client

func initRedis() {
	redisDatabase = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if _, err := redisDatabase.Ping(ctx).Result(); err != nil {
		log.Fatal("!@#$% Redis connection failed: ", err)
	} else {
		log.Println("Connected to Redis")
	}
}

func saveKeyVal(key, val string) {
	if err := redisDatabase.Set(ctx, key, val, 0).Err(); err != nil {
		log.Println("!@#$% Redis set error: ", err)
	} else {
		log.Printf("Saved %s:%s to redis\n", key, val)
	}
}
