package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	configPath := os.Getenv("CONFIG_PATH")
	logrus.Println("Config Path", configPath)
	if configPath == "" {
		panic("CONFIG_PATH not set")
	}
	viper.SetConfigName("app")
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("panicking; error reading config file: %v, configPath: %s;", err, configPath))
	}
}
