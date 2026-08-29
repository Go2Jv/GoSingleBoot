package main

import (
	"GoSingleBoot/internal/config"
	"GoSingleBoot/internal/db"
	"GoSingleBoot/internal/logger"
	"GoSingleBoot/internal/router"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config.InitConfig()

	logger.InitLogger()
	defer logger.Logger.Sync()

	db.NewPostgresClient()
	defer db.Client.Close()

	startHttpServer()
}

func startHttpServer() {

	mainRouter := router.MainRouter()
	server := http.Server{
		Addr:         config.Config.ApplicationCfg.Port,
		Handler:      mainRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		logger.Logger.Info("后端服务启动成功,在 " + config.Config.ApplicationCfg.Port + " 运行")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("后端服务启动失败: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Info("正在关闭后端服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Logger.Error("后端服务关闭失败 :" + err.Error())
	}

}
