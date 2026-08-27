package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Cfg struct {
	ApplicationCfg ApplicationCfg `json:"Application"`
	DatabaseCfg    DatabaseCfg    `json:"Database"`
	JwtCfg         JwtCfg         `json:"Jwt"`
	CorsCfg        CorsCfg        `json:"Cors"`
}

var Config Cfg

type ApplicationCfg struct {
	Name string `json:"Name"`
	Port string `json:"Port"`
	Text bool   `json:"Text"`
}

type DatabaseCfg struct {
	Master string   `json:"Master"`
	Slaves []string `json:"Slaves"`
}

type JwtCfg struct {
	StringSecretKey string `json:"SecretKey"`
	SecretKey       []byte `json:"-"`
	Expiration      int64  `json:"Expiration"`
	Issuer          string `json:"Issuer"`
}

type CorsCfg struct {
	AllowedOrigins     []string `json:"AllowedOrigins"`
	AllowedCredentials bool     `json:"AllowedCredentials"`
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

	// 这里直接转换！方便我直接调用！
	cfg.JwtCfg.SecretKey = []byte(cfg.JwtCfg.StringSecretKey)

	if cfg.CorsCfg.AllowedCredentials && len(cfg.CorsCfg.AllowedOrigins) == 0 {
		panic("你设置了AllowedCredentials 为true，那么不可以设置AllowedOrigins为空 !")
	}

	for _, origin := range cfg.CorsCfg.AllowedOrigins {
		if origin == "*" && cfg.CorsCfg.AllowedCredentials {
			panic("你设置了AllowedCredentials 为true，那么不可以设置AllowedOrigins为* !")
		}
	}

	if !cfg.CorsCfg.AllowedCredentials && len(cfg.CorsCfg.AllowedOrigins) == 0 {
		cfg.CorsCfg.AllowedOrigins = []string{"*"}
	}
	fmt.Println(cfg)
	Config = cfg
}
