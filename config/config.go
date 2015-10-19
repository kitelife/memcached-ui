package config

import (
	"encoding/json"
	"io/ioutil"
)

type ServerConfig struct {
	Source          string
	MiddlemanName   string
	MiddlemanConfig map[string]string
}

type MUConfigStruct struct {
	Servers map[string]ServerConfig
}

func LoadMUConfig(configPath string) (conf MUConfigStruct, err error) {
	confContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return conf, err
	}
	err = json.Unmarshal(confContent, &conf)
	return
}
