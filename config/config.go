package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	BotName      string `json:"bot_name"`
	BotUser      string `json:"user"`
	BotPass      string `json:"pass"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Port         int    `json:"port"`
	Subreddit    string `json:"subreddit"`
}

func LoadConfig(configPath string) Config {
	_, err := os.Stat(configPath)
	if err != nil {
		finalErr := fmt.Errorf("Config file %s not found. Reason: %s", configPath, err.Error())
		panic(finalErr)
	}

	buff, err := ioutil.ReadFile(configPath)
	if err != nil {
		finalErr := fmt.Errorf("Could not read config file %s. Reason: %s", configPath, err.Error())
		panic(finalErr)
	}

	conf := new(Config)
	err = json.Unmarshal(buff, conf)
	if err != nil {
		finalErr := fmt.Errorf("Failed to parse Configuration file %s. Reason: %s", configPath, err.Error())
		panic(finalErr)
	}

	return *conf
}
