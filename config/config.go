package config

import (
	"encoding/json"
	"io/ioutil"
)

type YiiConfigStruct struct {
	Status  string
	AppName string
	Hash    string
}

type ServerConfig struct {
	Alias string
	Yii   YiiConfigStruct
}

type AppConfigStruct struct {
	Servers map[string]ServerConfig
}

func LoadAppConfig(configPath string) (conf AppConfigStruct, err error) {
	confContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return conf, err
	}
	err = json.Unmarshal(confContent, &conf)
	return
}
