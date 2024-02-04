package storage

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

var redisClient *redis.Client

func RedisInit() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "", // no password set
		DB:           3,  // use default DB
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	})
	err := redisClient.Ping().Err()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("redis connected!")
}
