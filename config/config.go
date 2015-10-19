package config

import (
	"encoding/json"
	"io/ioutil"
)

type ServerConfig struct {
	Alias           string
	MiddlemanName   string
	MiddlemanConfig map[string]string
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
