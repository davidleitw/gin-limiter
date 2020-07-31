package limiter

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func LimitMiddle(lc *LimitController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now().Unix()
		path := ctx.FullPath()
		method := ctx.Request.Method
		ipAddress := ctx.ClientIP()

		globalIndex := "Source:" + ipAddress     // Source:123.456.78.9
		singleIndex := path + method + ipAddress // /a/postpost/123.456.78.9
		limit := lc.GetSingleLimit(path, method)
		globalLimit := lc.GetGlobalLimit()

		fmt.Println(ipAddress, now, globalIndex, singleIndex, limit, globalLimit)
		// good request
		ctx.Next()
	}
}
