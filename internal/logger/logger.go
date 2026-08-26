package logger

import (
	"GoSingleBoot/internal/config"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
	text := config.Config.ApplicationCfg.Text
	if text {
		logger, err := zap.NewProduction()
		if err != nil {
			panic("初始化Logger 时发生错误 : " + err.Error())
		}
		Logger = logger
	} else {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic("初始化Logger 时发生错误 : " + err.Error())
		}
		Logger = logger
	}
	Logger.Info("日志功能初始化成功")
}
