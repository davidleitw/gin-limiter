package limiter

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func LimitMiddle(lc *LimitController) gin.HandlerFunc {
	lc.Init()

	return func(ctx *gin.Context) {
		now := time.Now().Unix()
		path := ctx.FullPath()
		method := ctx.Request.Method
		ipAddress := ctx.ClientIP()

		globalKey := "Source:" + ipAddress // Source:123.456.78.9
		globalLimit := lc.GetGlobalLimit()

		singleKey := path + "/" + method + "/" + ipAddress // /a/post/post/123.456.78.9
		singleLimit := lc.GetSingleLimit(path, method)

		script := redis.NewScript(Script)
		args := []interface{}{now, globalLimit, singleLimit}
		keys := []string{globalKey, singleKey}

		result, err := script.Run(ctx, lc.RedisDB, keys, args).Result()
		if err != nil {
			fmt.Println("Script run error = ", err)
		}
		fmt.Println("result = ", result)

		ctx.Next()
	}
}
