package main

import (
	"context"
	"log"

	limiter "github.com/davidleitw/gin-limiter"
	"github.com/go-redis/redis/v8"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	arg1 := false
	sHA, _ := rdb.ScriptLoad(context.Background(), limiter.TestScript).Result()
	// 		result := lc.RedisDB.EvalSha(context.Background(), lc.GetShaScript(), keys, args)

	for i := 0; i < 5; i++ {
		result := rdb.EvalSha(context.Background(), sHA, []string{}, []interface{}{arg1})
		log.Println("result = ", result.Val())
	}

}
