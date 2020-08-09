package limiter

import (
	"context"
	"net/http"
	"strconv"
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

			singleKey := path + "/" + method + ":" + ipAddress // /a/post/POST:123.456.78.9
			singleLimit := lc.GetSingleLimit(path, method)     // Single router 限制次數

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

			results, err := lc.RedisDB.EvalSha(context.Background(), lc.GetShaScript(), keys, args).Result()
			if err != nil {
				// code 500, server error
				ctx.JSON(http.StatusInternalServerError, err)
				ctx.Abort()
			}

			result := results.([]interface{})
			globalRemaining := result[0].(int64) // global 剩餘次數
			singleRemaining := result[1].(int64) // single router 剩餘次數

			lc.logger.Printf("|%s|  %s  \"%s\" | global remain=%d | single remain=%d | DeadLine{%s/%s}\n", ipAddress, method, path,
				globalRemaining, singleRemaining, lc.globalRate.GetDeadLineFormat(), lc.routerRates.GetDeadLineFormat(path, method))

			if globalRemaining == -1 {
				ctx.JSON(http.StatusTooManyRequests, "To many request!")
				ctx.Header("X-RateLimit-Reset-global", lc.globalRate.GetDeadLineFormat())
				ctx.Abort()
			}

			if singleRemaining == -1 {
				ctx.JSON(http.StatusTooManyRequests, "To many request!")
				ctx.Header("X-RateLimit-Reset-single", lc.routerRates.GetDeadLineFormat(path, method))
				ctx.Abort()
			}

			ctx.Header("X-RateLimit-Limit-global", strconv.Itoa(globalLimit))
			ctx.Header("X-RateLimit-Remaining-global", strconv.FormatInt(globalRemaining, 10))
			ctx.Header("X-RateLimit-Reset-global", lc.globalRate.GetDeadLineFormat())
			ctx.Header("X-RateLimit-Limit-single", strconv.Itoa(singleLimit))
			ctx.Header("X-RateLimit-Remaining-single", strconv.FormatInt(singleRemaining, 10))
			ctx.Header("X-RateLimit-Reset-single", lc.routerRates.GetDeadLineFormat(path, method))

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

			singleKey := path + "/" + method + ":" + ipAddress // /a/post/post:123.456.78.9
			singleLimit := lc.GetSingleLimit(path, method)     // Single router 限制次數

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

			results, err := lc.RedisDB.EvalSha(context.Background(), lc.GetShaScript(), keys, args).Result()
			if err != nil {
				// code 500, server error
				ctx.JSON(http.StatusInternalServerError, err)
				ctx.Abort()
			}

			result := results.([]interface{})
			globalRemaining := result[0].(int64) // global 剩餘次數
			singleRemaining := result[1].(int64) // single router 剩餘次數

			if globalRemaining == -1 {
				ctx.JSON(http.StatusTooManyRequests, "To many request!")
				ctx.Header("X-RateLimit-Reset-global", lc.globalRate.GetDeadLineFormat())
				ctx.Abort()
			}

			if singleRemaining == -1 {
				ctx.JSON(http.StatusTooManyRequests, "To many request!")
				ctx.Header("X-RateLimit-Reset-single", lc.routerRates.GetDeadLineFormat(path, method))
				ctx.Abort()
			}

			ctx.Header("X-RateLimit-Limit-global", strconv.Itoa(globalLimit))
			ctx.Header("X-RateLimit-Remaining-global", strconv.FormatInt(globalRemaining, 10))
			ctx.Header("X-RateLimit-Reset-global", lc.globalRate.GetDeadLineFormat())
			ctx.Header("X-RateLimit-Limit-single", strconv.Itoa(singleLimit))
			ctx.Header("X-RateLimit-Remaining-single", strconv.FormatInt(singleRemaining, 10))
			ctx.Header("X-RateLimit-Reset-single", lc.routerRates.GetDeadLineFormat(path, method))

			ctx.Next()
		}
	}
}
