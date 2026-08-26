package router

import "github.com/gin-gonic/gin"

func MainRouter() *gin.Engine {
	mainRouter := gin.New()
	return mainRouter
}
