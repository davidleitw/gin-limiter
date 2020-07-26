package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func exampleNewClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	fmt.Println(pong, err)

	err = rdb.Set(context.Background(), "davidleitw", "123", 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	err = rdb.Get(context.Background(), "davidlettw").Err()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	exampleNewClient()
}
