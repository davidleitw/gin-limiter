package main

import (
	"context"
	"fmt"

	limiter "github.com/davidleitw/gin-limiter"

	"github.com/go-redis/redis/v8"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	sha, _ := rdb.ScriptLoad(context.Background(), limiter.TestScript).Result()
	result, _ := rdb.EvalSha(context.Background(), sha, []string{}).Result()
	fmt.Println(result)
}
