package router

import (
	"GoSingleBoot/internal/config"
	"GoSingleBoot/internal/docs"
	"GoSingleBoot/internal/middleware"
	"net/http"

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

	// 添加了openapi自动生成和一个swagger页面
	text := config.Config.ApplicationCfg.Text
	if text {
		docs.GenerateOpenAPI()
		mainRouter.GET("/docs", docs.Handler)
	}

	// Register Router 注册路由
	rg := mainRouter.Group("/api")
	RegisterLoginRouter(rg)

	// NoRoute 404
	mainRouter.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "不存在该Api",
		})
	})
	return mainRouter
}
