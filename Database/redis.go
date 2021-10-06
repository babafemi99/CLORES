package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var Cache *redis.Client
var CacheChannel chan string

func SetupRedis() {
	Cache = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
}

func SetupCacheChannel() {

	go func(ch chan string) {
		for {

			time.Sleep(5 * time.Second)
			Cache.Del(context.Background(), <-ch)
			fmt.Print("cache clear")
		}
	}(CacheChannel)
}

func ClearCache(keys ...string) {
	for _,key := range keys {
		CacheChannel <- key
	}
}
