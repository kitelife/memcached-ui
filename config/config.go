package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type InstanceConfig struct {
	Source          string            `yaml:"source"`
	MiddlemanName   string            `yaml:"middelman_name"`
	MiddlemanConfig map[string]string `yaml:"middleman_config"`
}

type BasicAuthConfig struct {
	On       bool   `yaml:"on"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type AppConfigStruct struct {
	Instances  map[string]InstanceConfig `yaml:"instances"`
	Basic_auth BasicAuthConfig           `yaml:"basic_auth"`
}

var DefaultAppConfig = AppConfigStruct{
	Instances: map[string]InstanceConfig{
		"localhost-Yii": InstanceConfig{
			Source:        "localhost:11211",
			MiddlemanName: "yii",
			MiddlemanConfig: map[string]string{
				"appName":            "gxt",
				"hash":               "yes",
				"php_bin":            "php",
				"unserialize_script": "./middleman/middleman/unserialize_to_json.php",
			},
		},
	},
	Basic_auth: BasicAuthConfig{
		On:       true,
		Username: "test",
		Password: "test",
	},
}

func LoadAppConfig(configPath string) (AppConfigStruct, error) {
	confContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return DefaultAppConfig, err
	}
	var conf AppConfigStruct
	if err := yaml.Unmarshal(confContent, &conf); err != nil {
		return DefaultAppConfig, err
	}
	return conf, nil
}
