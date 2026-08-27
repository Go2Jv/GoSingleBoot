package middleware

import (
	"GoSingleBoot/internal/bizErr"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GlobalErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if len(c.Errors) != 0 {
				err := c.Errors.Last().Err

				// 90%都是业务错误
				var bz *bizErr.BizErr
				if ok := errors.As(err, &bz); ok {
					c.JSON(http.StatusOK, gin.H{
						"code": bz.Code,
						"msg":  bz.Msg,
					})
					return
				}

				// 兜底的，未知错误
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  err.Error(),
				})
			}
		}()
		c.Next()
	}
}
