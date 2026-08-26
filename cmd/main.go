package main

import (
	"GoSingleBoot/internal/config"
	"GoSingleBoot/internal/db"
	"GoSingleBoot/internal/logger"
)

func main() {
	config.InitConfig()
	logger.InitLogger()
	defer logger.Logger.Sync()
	db.NewPostgresClient()
	defer db.Client.Close()

}
