package middleware

import (
	"GoSingleBoot/internal/bizErr"
	"GoSingleBoot/internal/logger"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GlobalErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if len(c.Errors) != 0 {
				err := c.Errors.Last().Err

				// 99%都是业务错误，正常使用向c.Error(bz)必须是*BizErr类型的！！！
				var bz *bizErr.BizErr
				if ok := errors.As(err, &bz); ok {
					// logger打的信息和返回给用户的信息是不一样的！
					logger.Logger.Warn(bz.Log.Error())
					c.JSON(http.StatusOK, gin.H{
						"code": bz.Code,
						"msg":  bz.Msg,
					})
					return
				}

				// 兜底的，主要是防止有SB或是实习生直接向c.Error()直接传入，没有使用bizErr.Wrap()封装
				// 未知错误
				// 正常使用向c.Error(bz)必须是*BizErr类型的！！！
				logger.Logger.Warn(err.Error())
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "服务器繁忙，请稍后重试",
				})
			}

		}()
		c.Next()
	}
}
