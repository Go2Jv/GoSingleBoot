package bizErr

import "github.com/gin-gonic/gin"

// Abort 在if err != nil {}中使用，要手动return.
// 作用只有手动将错误传到 全局错误处理中间件中
func Abort(c *gin.Context, code int, msg string, err error) {
	bz := Wrap(code, msg, err)
	_ = c.Error(bz)
}
