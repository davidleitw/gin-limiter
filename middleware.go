package limiter

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func (lc *LimitController) GenerateLimitMiddleWare() gin.HandlerFunc {
	lc.Init()

	if lc.Mode() == "debug" {
		return func(ctx *gin.Context) {
			var globalExpired bool = false
			var singleExpired bool = false

			now := time.Now().Unix()
			path := ctx.FullPath()
			method := ctx.Request.Method
			ipAddress := ctx.ClientIP()

			globalKey := "Source:" + ipAddress // Source:123.456.78.9
			globalLimit := lc.GetGlobalLimit() // Global 限制次數
			// globalDeadLine := lc.globalRate.GetDeadLine()

			singleKey := path + "/" + method + ":" + ipAddress // /a/post/post:123.456.78.9
			singleLimit := lc.GetSingleLimit(path, method)     // Single router 限制次數
			// singleDeadLine := lc.routerRates.GetDeadLine(path, method)

			if now > lc.globalRate.GetDeadLine() {
				globalExpired = true
				lc.globalRate.UpdateDeadLine() // 更新DeadLine
			}
			if now > lc.routerRates.GetDeadLine(path, method) {
				singleExpired = true
				lc.routerRates.UpdateDeadLine(path, method) // 更新單一path DeadLine
			}

			keys := []string{globalKey, singleKey}
			args := []interface{}{globalLimit, singleLimit, globalExpired, singleExpired}

			result := lc.RedisDB.EvalSha(context.Background(), lc.GetShaScript(), keys, args)

			log.Printf("now: %d, global deadline = %d, single router deadline = %d\n", now, lc.globalRate.GetDeadLine(), lc.routerRates.GetDeadLine(path, method))
			log.Printf("global expired = %t, single expired = %t\n", globalExpired, singleExpired)
			log.Printf("Request Information: global{Key:%s, Limit:%d} single{Key:%s, Limit:%d}\n", globalKey, globalLimit, singleKey, singleLimit)
			log.Println("result = ", result.Val())

			ctx.Next()
		}
	} else {
		return func(ctx *gin.Context) {
			var globalExpired bool = false
			var singleExpired bool = false

			now := time.Now().Unix()
			path := ctx.FullPath()
			method := ctx.Request.Method
			ipAddress := ctx.ClientIP()

			globalKey := "Source:" + ipAddress // Source:123.456.78.9
			globalLimit := lc.GetGlobalLimit() // Global 限制次數
			// globalDeadLine := lc.globalRate.GetDeadLine()

			singleKey := path + "/" + method + ":" + ipAddress // /a/post/post:123.456.78.9
			singleLimit := lc.GetSingleLimit(path, method)     // Single router 限制次數
			// singleDeadLine := lc.routerRates.GetDeadLine(path, method)

			if now > lc.globalRate.GetDeadLine() {
				globalExpired = true
				lc.globalRate.UpdateDeadLine() // 更新DeadLine
			}
			if now > lc.routerRates.GetDeadLine(path, method) {
				singleExpired = true
				lc.routerRates.UpdateDeadLine(path, method) // 更新單一path DeadLine
			}

			keys := []string{globalKey, singleKey}
			args := []interface{}{globalLimit, singleLimit, globalExpired, singleExpired}

			result := lc.RedisDB.EvalSha(context.Background(), lc.GetShaScript(), keys, args)
			log.Println(result.Val())

			ctx.Next()
		}
	}
}
