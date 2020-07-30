package limiter

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func LimitMiddle(lc *LimitController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.FullPath()
		fmt.Println(path)

		// good request
		ctx.Next()
	}
}
