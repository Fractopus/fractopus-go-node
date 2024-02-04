package storage

import (
	"github.com/go-redis/redis"
	"github.com/tidwall/gjson"
	"log"
	"time"
)

const defaultExpiry = 60 * time.Minute
const keyPre = "fractopus:"

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

func SetWithDefaultExpiry(key string, value interface{}) error {
	return SetWithTimeout(key, value, defaultExpiry)
}

func SetWithTimeout(key string, value interface{}, expiry time.Duration) error {
	return redisClient.Set(keyPre+key, value, expiry).Err()
}

func GetJsonValue(key string) (gjson.Result, error) {
	result, err := redisClient.Get(keyPre + key).Result()
	if err != nil {
		return gjson.Result{}, err
	}
	return gjson.Parse(result), nil
}
func GetStringValue(key string) (string, error) {
	result, err := redisClient.Get(keyPre + key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
