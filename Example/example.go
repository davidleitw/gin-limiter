package main

import (
	"context"
	"fmt"
	"log"

	limiter "github.com/davidleitw/gin-limiter"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func Test() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, _ := rdb.Ping(context.Background()).Result()
	fmt.Println(pong)
}

func NewServer() *gin.Engine {
	server := gin.Default()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	limitControl, _ := limiter.DefaultController(rdb, "24-H", 21000)

	server.POST("/post1", post1) // /post1
	_ = limitControl.Add("/post", "20-M", "post", 120)

	server.POST("api/post2", post2) // /api/post2
	_ = limitControl.Add("api/post2", "15-H", "post", 200)
	return server
}

func post1(ctx *gin.Context) {
	ctx.String(200, ctx.FullPath())
}

func post2(ctx *gin.Context) {
	ctx.String(200, ctx.FullPath())
}

func main() {
	server := NewServer()
	err := server.Run(":8080")
	if err != nil {
		log.Println(err)
	}
	// Test()
}
