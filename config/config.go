package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadEnv(path, filename, ext string) {
	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType(ext)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config file: %s", err.Error())
	}
}
