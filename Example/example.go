package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
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

	val, err := rdb.Get(Context, "davidleitw2").Int()
	if err != nil {
		fmt.Println("err = ", err)
	}
	fmt.Println("value = ", val)

	val++

	err = rdb.Set(Context, "davidleitw2", val, 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	val, err = rdb.Get(Context, "davidleitw2").Int()
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

func NewServer() *gin.Engine {
	server := gin.Default()

	return server
}

func main() {
	exampleNewClient()
	t0 := time.Now()
	t1 := time.Now()
	//设置期间经历了50秒时间
	t2 := time.Now().Add(time.Second * 50)
	fmt.Println(t2.Sub(t1)) //t2与t1相差： 50s
	fmt.Println(t1.Sub(t2))
	fmt.Println(t1.Before(t2))
	fmt.Println(t1.Before(t0))
	fmt.Println(t0, t1)
	// server := NewServer()
	// server.Run()

	// l := limiter.LimitController()
}
