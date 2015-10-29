package config

import (
	"encoding/json"
	"io/ioutil"
)

type InstanceConfig struct {
	Source          string
	MiddlemanName   string
	MiddlemanConfig map[string]string
}

type BasicAuthConfig struct {
	On       string
	Username string
	Password string
}

type AppConfigStruct struct {
	Instances  map[string]InstanceConfig
	Basic_auth BasicAuthConfig
}

func LoadAppConfig(configPath string) (conf AppConfigStruct, err error) {
	confContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return conf, err
	}
	err = json.Unmarshal(confContent, &conf)
	return
}
