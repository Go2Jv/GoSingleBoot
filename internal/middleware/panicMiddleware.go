package middleware

import (
	"GoSingleBoot/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PanicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				e, ok := err.(error)
				if ok {
					logger.Logger.Error("发生 panic", zap.Error(e))
				} else {
					logger.Logger.Error("发生 panic", zap.Any("panic", err))
				}
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "服务器繁忙请稍后尝试",
				})
				return
			}
		}()
		c.Next()
	}
}
