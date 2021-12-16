package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig(path string) {
	config := path
	viper.SetConfigFile(config)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}
}

func ReadConfig(key string) string {
	viper.ReadInConfig()
	return viper.GetString(key)
}
