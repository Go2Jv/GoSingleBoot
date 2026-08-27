package router

import (
	"GoSingleBoot/internal/loginHandler"

	"github.com/gin-gonic/gin"
)

func RegisterLoginRouter(rg *gin.RouterGroup) {
	//实例化LoginHandler ->Login
	lhandler := loginHandler.NewLoginHandler()

	rg.POST("/login", lhandler.Login)
}
