package bizErr

import (
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
)

// 这个文件用来封装一些常用的方法 , 如果err != nil -> false
// false 是有问题就ok ，在我这里
// true 是没有问题！

// Check 检测err是否为nil ,如果是让全局中间件返回code和msg ,让全局中间件的logger打印err
func Check(c *gin.Context, code int, msg string, err error) bool {
	if err != nil {
		bz := Wrap(code, msg, err)
		_ = c.Error(bz)
		return false
	}
	return true
}

// ServerBusy 直接封装了很多不方便的错误，如生成token失败，或是读取文件失败，数据库/Redis崩了
// 你总不可以直接和用户将服务器文件读取失败吧？？？或是我的数据库崩了吧？？？
// 直接返回服务器繁忙骗一下用户
func ServerBusy(c *gin.Context, err error) bool {
	if err != nil {
		bz := Wrap(500, "服务器繁忙，请稍后重试", err)
		_ = c.Error(bz)
		return false
	}
	return true
}

// Validation 用来验证参数时报错时
func Validation(c *gin.Context, msg string, err error) bool {
	if msg == "" {
		msg = "非法传参"
	}
	if ok := Check(c, 400, msg, err); !ok {
		return false
	}
	return true
}

// SQLNotFound 查询不到数据->Limit(1)
func SQLNotFound(c *gin.Context, msg string, err error) bool {
	if msg == "" {
		msg = "不存在该数据"
	}
	if errors.Is(err, sql.ErrNoRows) {
		_ = c.Error(Wrap(404, msg, err))
		return false
	}
	if err != nil {
		_ = c.Error(Wrap(500, "服务器繁忙，请稍后重试", err))
		return false
	}

	return true
}
