package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var CfgPath = "/home/leo/go/src/telegram_bot/config/config.yaml"

type Config struct {
	TelegramAPIToken string `yaml:"BOT_TOKEN"`
	DbUrl            string `yaml:"DB_URL"`
}

func LoadConfig() (*Config, error) {
	buf, err := ioutil.ReadFile(CfgPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
