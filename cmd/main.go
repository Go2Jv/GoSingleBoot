package main

import (
	"GoSingleBoot/internal/config"
	"GoSingleBoot/internal/db"
)

func main() {
	config.InitConfig()
	db.NewPostgresClient()

}
