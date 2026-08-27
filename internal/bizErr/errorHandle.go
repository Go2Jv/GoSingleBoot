package bizErr

import (
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
)

// 我这里只是负责抛出或是传递错误，错误应该让错误中间件处理，打日志也是，防止重复打日志！

func Wrap(code int, msg string, err error) *BizErr {
	// 如果err为空，直接使用msg！方便快速开发
	if err == nil {
		err = errors.New(msg)
	}
	return &BizErr{
		Code: code,
		Msg:  msg,
		Log:  err,
	}
}

// Throw 自动打包为Wrap()，然后抛出到全局错误处理中间件
func Throw(c *gin.Context, code int, msg string, err error) {
	// 如果err为空，直接使用msg！方便快速开发
	if err == nil {
		err = errors.New(msg)
	}
	bz := Wrap(code, msg, err)
	_ = c.Error(bz)
}

// Validation 用来验证参数时报错时
func Validation(c *gin.Context, err error) bool {
	if err != nil {
		var bz = &BizErr{
			Code: 400,
			Msg:  "非法传参",
			Log:  err,
		}
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
			Log:  err,
		}
		_ = c.Error(bz)
		return false
	}
	if err != nil {
		panic("数据库异常 :" + err.Error())
	}
	return true
}
