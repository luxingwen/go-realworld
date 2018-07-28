package config

import (
	"github.com/go-ini/ini"
	"log"
)

var (
	ServerConf = &ServerConfig{}
	MySqlConf  = &MySqlConfig{}
)

type ServerConfig struct {
	Port     string
	RunMode  string
	FilePath string
}

type MySqlConfig struct {
	Host     string
	Port     string
	DbName   string
	UserName string
	PassWd   string
}

func init() {
	cfg, err := ini.Load("conf/app.conf")
	if err != nil {
		log.Fatal("load app conf err:%v", err)
	}
	err = cfg.Section("server").MapTo(ServerConf)
	if err != nil {
		log.Fatal("init server conf err:", err)
	}
	err = cfg.Section("mysql").MapTo(MySqlConf)
	if err != nil {
		log.Fatal("init mysql conf err:", err)
	}
}
