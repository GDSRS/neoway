package config

import (
	"github.com/spf13/viper"
)

func LoadConfig(configPath string) {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	}

	viper.ReadInConfig()
}
