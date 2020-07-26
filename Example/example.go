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

	Context := context.Background()

	pong, err := rdb.Ping(Context).Result()
	fmt.Println(pong, err)

	err = rdb.Set(Context, "davidleitw2", "123", 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	val, err := rdb.Get(Context, "davidleitw2").Result()
	if err != nil {
		fmt.Println("err = ", err)
	}
	fmt.Println("value = ", val)

	val2, err := rdb.Get(Context, "NotExistValue").Result()
	if err == redis.Nil {
		fmt.Println("val2 not exist")
	} else if err != nil {
		fmt.Println("err = ", err)
	} else {
		fmt.Println("value2 = ", val2)
	}
}

func main() {
	exampleNewClient()
}
