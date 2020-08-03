package limiter

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LimitMiddle(lc *LimitController) gin.HandlerFunc {
	lc.Init()

	return func(ctx *gin.Context) {
		var globalExpired string = "false"
		var singleExpired string = "false"

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
			globalExpired = "true"
			lc.globalRate.UpdateDeadLine() // 更新DeadLine
		}
		if now > lc.routerRates.GetDeadLine(path, method) {
			singleExpired = "true"
			lc.routerRates.UpdateDeadLine(path, method) // 更新單一path DeadLine
		}

		args := []interface{}{now, globalLimit, singleLimit, globalExpired, singleExpired}
		keys := []string{globalKey, singleKey}

		result := lc.RedisDB.EvalSha(context.Background(), lc.GetShaScript(), keys, args)
		log.Println(result)

		ctx.Next()
	}
}
