package main

import (
	"log"

	"net/http"

	limiter "github.com/davidleitw/gin-limiter"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	server := gin.Default()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})

	dispatcher, err := limiter.LimitDispatcher("24-M", 100, rdb)
	if err != nil {
		log.Println(err)
	}

	server.POST("/ExamplePost1", dispatcher.MiddleWare("4-M", 20), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello ExamplePost1")
	})

	server.GET("/ExampleGet1", dispatcher.MiddleWare("5-M", 10), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello ExampleGet1")
	})

	err = server.Run(":8080")
	if err != nil {
		log.Println("gin server error = ", err)
	}
}
