package handler

import (
	"GoSingleBoot/internal/bizErr"
	"GoSingleBoot/internal/db"
	"GoSingleBoot/internal/dto/req"
	"GoSingleBoot/internal/dto/resp"
	"GoSingleBoot/internal/jwt"
	"GoSingleBoot/internal/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (lc *LoginHandler) Login(c *gin.Context) {
	var rq req.LoginByPassword
	err := c.ShouldBindJSON(&rq)
	// 原来要写一大堆的逻辑...
	//if err != nil {
	//	logger.Warn(err.Error())
	//	bizErr.Validation(c, err)
	//	return
	//}
	if ok := bizErr.Validation(c, "", err); !ok {
		return
	}

	var user model.User
	err = db.Client.Master.NewSelect().
		Model(&user).
		Where("username = ?", rq.Username).
		Limit(1).
		Scan(c, &user)
	//if err != nil {
	//	if err == sql.ErrNoRows {
	//		xxxx 一大堆逻辑
	//	}
	//}
	if ok := bizErr.SQLNotFound(c, "账号为注册", err); !ok {
		return
	}

	if user.Password != rq.Password {
		bizErr.Abort(c, 400, "账号或密码错误", nil)
		return
	}

	token, err := jwt.GenerateToken(strconv.Itoa(int(user.ID)))
	if err != nil {
		bizErr.Abort(c, 500, err.Error(), nil)
		return
	}

	var rs resp.CodeMsgAndData
	rs.Code = 200
	rs.Msg = "登陆成功"
	rs.Data = token
	logger.Info(fmt.Sprintf("用户id : %v 用户: %v 登陆成功,token为 %v ", user.ID, user.Username, token))
	c.JSON(http.StatusOK, rs)

}
