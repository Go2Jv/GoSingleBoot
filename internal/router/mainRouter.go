package router

import (
	"GoSingleBoot/internal/config"
	"GoSingleBoot/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setMode() {
	text := config.Config.ApplicationCfg.Text
	if !text {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func MainRouter() *gin.Engine {
	setMode()
	mainRouter := gin.New()
	// 中间件
	mainRouter.Use(middleware.CorsMiddleware())
	mainRouter.Use(middleware.PanicMiddleware())
	mainRouter.Use(middleware.GlobalErrorMiddleware())

	// Register Router 注册路由
	rg := mainRouter.Group("/api")
	RegisterLoginRouter(rg)

	return mainRouter
}
