package bizErr

import (
	"GoSingleBoot/internal/logger"
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Warp(code int, msg string) *BizErr {
	return &BizErr{
		Code: code,
		Msg:  msg,
	}
}

// Throw 自动打包为Warp()，然后抛出到全局错误处理中间件
func Throw(c *gin.Context, code int, msg string) {
	bz := Warp(code, msg)
	_ = c.Error(bz)
}

// Validation 用来验证参数时报错时
func Validation(c *gin.Context, err error) bool {
	if err != nil {
		var bz = &BizErr{
			Code: 400,
			Msg:  "非法传参",
		}
		logger.Logger.Warn(err.Error())
		_ = c.Error(bz)
		return false
	}
	return true

}

// SQLNotFound 查询不到数据->Limit(1)
func SQLNotFound(c *gin.Context, err error) bool {
	if errors.Is(err, sql.ErrNoRows) {
		var bz = &BizErr{
			Code: 404,
			Msg:  "不存在",
		}
		logger.Logger.Warn(err.Error())
		_ = c.Error(bz)
		return false
	}
	if err != nil {
		logger.Logger.Error(err.Error(), zap.Stack("stack"))
		panic("数据库异常 :" + err.Error())
	}
	return true
}
