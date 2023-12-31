package config

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/scorpiotzh/toolib"
)

var (
	Cfg CfgServer
	log = logger.NewLogger("config", logger.LevelDebug)
)

func InitCfg(configFilePath string) error {
	if configFilePath == "" {
		configFilePath = "./config/config.yaml"
	}
	log.Debug("config file path：", configFilePath)
	if err := toolib.UnmarshalYamlFile(configFilePath, &Cfg); err != nil {
		return fmt.Errorf("UnmarshalYamlFile err:%s", err.Error())
	}
	log.Debug("config file：", toolib.JsonString(Cfg))
	return nil
}

func AddCfgFileWatcher(configFilePath string) (*fsnotify.Watcher, error) {
	if configFilePath == "" {
		configFilePath = "./config/config.yaml"
	}
	return toolib.AddFileWatcher(configFilePath, func() {
		log.Debug("config file path：", configFilePath)
		if err := toolib.UnmarshalYamlFile(configFilePath, &Cfg); err != nil {
			log.Error("UnmarshalYamlFile err:", err.Error())
		}
		log.Debug("config file：", toolib.JsonString(Cfg))
	})
}

type CfgServer struct {
	Server struct {
		key      string `json:"key" yaml:"key"`
		HttpAddr string `json:"http_addr" yaml:"http_addr"`
	} `json:"server" yaml:"server"`
	IpWhitelist map[string]string `json:"ip_whitelist" yaml:"ip_whitelist"`
	Notify      struct {
		LarkErrorKey string `json:"lark_error_key" yaml:"lark_error_key"`
		SentryDsn    string `json:"sentry_dsn" yaml:"sentry_dsn"`
	} `json:"notify" yaml:"notify"`
	DB struct {
		Mysql DbMysql `json:"mysql" yaml:"mysql"`
	} `json:"db" yaml:"db"`
}

type DbMysql struct {
	Addr     string `json:"addr" yaml:"addr"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	DbName   string `json:"db_name" yaml:"db_name"`
}

func (c *CfgServer) SetKey(key string) {
	c.Server.key = key
}

func (c *CfgServer) GetKey() string {
	return c.Server.key
}
