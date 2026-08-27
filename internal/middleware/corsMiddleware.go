package middleware

import (
	"GoSingleBoot/internal/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 允许访问的域名列表（也可以用 AllowAllOrigins: true）
		AllowOrigins:     config.Config.CorsCfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: config.Config.CorsCfg.AllowedCredentials,
		// 预检请求 (OPTIONS) 的缓存时间，避免频繁发 OPTIONS 请求
		MaxAge: 3 * 30 * 24 * time.Hour,
	})
}
