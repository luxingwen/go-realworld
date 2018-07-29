package config

import (
	"log"

	"github.com/go-ini/ini"
)

var (
	ServerConf = &ServerConfig{}
	MySqlConf  = &MySqlConfig{}
	EmaiConf   = &EMailConfig{}
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

type EMailConfig struct {
	DefaultAdress string
	DefaultPort   int
	DefaultUser   string
	DefaultPasswd string
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
	err = cfg.Section("email").MapTo(EmaiConf)
	if err != nil {
		log.Fatal("init email conf err:", err)
	}
}
