package config

import (
	"encoding/json"
	"io/ioutil"
)

type AppConfigStruct struct {
	Servers map[string]string
}

func LoadAppConfig(configPath string) (conf AppConfigStruct, err error) {
	confContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return conf, err
	}
	err = json.Unmarshal(confContent, &conf)
	return
}
