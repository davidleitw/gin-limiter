package main

import (
	"log"

	limiter "github.com/davidleitw/gin-limiter"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func NewServer() *gin.Engine {
	server := gin.Default()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	limitControl, _ := limiter.DefaultController(rdb, "24-H", 2000)
	limitControl.Add("/post", "20-M", 120)

	server.POST("/post1", post1)    // /post1
	server.POST("api/post2", post2) // /api/post2
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
}
