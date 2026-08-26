package config

import (
	"encoding/json"
	"os"
)

type Cfg struct {
	ApplicationCfg ApplicationCfg
	OpenApi        bool `json:"OpenApi"`
	DatabaseCfg    DatabaseCfg
}

var Config Cfg

type ApplicationCfg struct {
	Name string `json:"Name"`
	Port string `json:"Port"`
}

type DatabaseCfg struct {
	Master string   `json:"Master"`
	Slaves []string `json:"Slaves"`
}

func InitConfig() {
	data, err := os.ReadFile("./config.json")
	if err != nil {
		panic("读取配置文件时报错 : " + err.Error())
	}
	var cfg Cfg
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		panic("解析Config.json文件时发生错误 : " + err.Error())
	}
	Config = cfg
}
