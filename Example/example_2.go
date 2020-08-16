package main

import (
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
	dispatcher, _ := limiter.DefaultController(rdb, "debug")
	diepatcher.Set("24-M", 100)

	server.POST("/post1", dispatcher.middle("20-S", 30), post1) // /post1

	server.POST("/api/post2", dispatcher.middle("15-M", 40), post2) // /api/post2

	server.POST("/post3", dispatcher.middle("11-D", 10), post3) // /post3

	return server
}

func post1(ctx *gin.Context) {
	ctx.String(200, ctx.FullPath())
}

func post2(ctx *gin.Context) {
	ctx.String(200, ctx.FullPath())
}

func post3(ctx *gin.Context) {
	ctx.String(200, ctx.ClientIP())
}

func main() {
	// server := NewServer()

	// err := server.Run(":8080")
	// if err != nil {
	// 	log.Println(err)
	// }
}
