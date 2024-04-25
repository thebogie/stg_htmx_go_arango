package utils

import (
	"github.com/spf13/viper"
	"sync"
)

var once sync.Once
var config *viper.Viper

func GetConfig() *viper.Viper {
	once.Do(func() {
		// Initialize Viper instance
		config = viper.New()

		// Set configuration options
		config.SetConfigFile(".env") // name of config file (without extension)

		// Read in the configuration file
		err := config.ReadInConfig()
		if err != nil {
			panic("Failed to read config file: " + err.Error())
		}
	})
	return config
}
