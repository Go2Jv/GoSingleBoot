package bizErr

import (
	"errors"
)

// 我这里只是负责抛出或是传递错误，错误应该让错误中间件处理，打日志也是，防止重复打日志！

// Wrap 只复杂包装，返回一个*BizErr
func Wrap(code int, msg string, err error) *BizErr {
	// 不可以都为空，我看哪一个sb干什么都不写的！！！
	if msg == "" && err == nil {
		panic("biz.Wrap传入的msg和err不可以同时为空！")
	}
	// 如果err为空，直接使用msg！方便快速开发
	if err == nil {
		err = errors.New(msg)
	}

	// 如果msg为空，直接使用err.Error()
	if msg == "" {
		msg = err.Error()
	}

	return &BizErr{
		Code: code,
		Msg:  msg,
		Log:  err,
	}
}
