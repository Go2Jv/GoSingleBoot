package router

import (
	"GoSingleBoot/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterLoginRouter(rg *gin.RouterGroup) {
	//实例化LoginHandler ->Login
	lhandler := handler.NewLoginHandler()

	rg.POST("/login", lhandler.Login)
}
