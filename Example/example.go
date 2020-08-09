package main

import (
	"log"

	limiter "github.com/davidleitw/gin-limiter"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	server := gin.Default()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})

	limitControl, err := limiter.DefaultController(rdb, "24-M", 100, "debug") // Debug mode, each 24 minutes can send 100 times request from single Ip.
	if err != nil {
		log.Println(err)
	}

	err = limitControl.Add("/ExamplePost1", "POST", "4-M", 20) // "/ExamplePost1" route, each 4 minutes can send 20 times request from single Ip.
	if err != nil {
		log.Println(err)
	}

	err = limitControl.Add("/ExampleGet1", "GET", "20-H", 40) // ".ExampleGet1" route, each 20 hours can send 40 times request from single Ip.
	if err != nil {
		log.Println(err)
	}
	server.Use(limitControl.GenerateLimitMiddleWare()) // create middleware

	server.POST("/ExamplePost1", func(ctx *gin.Context) {
		ctx.String(200, "Hello Example! In ExamplePost1")
	})

	server.GET("/ExampleGet1", func(ctx *gin.Context) {
		ctx.String(200, "Hello Example! In ExampleGet1")
	})

	err = server.Run(":8080")
	if err != nil {
		log.Println("gin server run error = ", err)
	}
}
