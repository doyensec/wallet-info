package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	IsDebug           bool   `yaml:"debug"`
	EtherscanEndpoint string `yaml:"etherscan_endpoint"`
	EtherscanApiKey   string `yaml:"etherscan_api_key"`
	SolscanEndpoint   string `yaml:"solscan_endpoint"`
}

var ActiveConfig *Config

func LoadConfig() error {
	file, err := os.Open("config.yml")
	if err != nil {
		return err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}

	ActiveConfig = &config
	return nil
}
