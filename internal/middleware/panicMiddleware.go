package middleware

import "github.com/gin-gonic/gin"

func PanicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
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
