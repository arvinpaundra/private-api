package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

func LoadEnv(path, filename, ext string) {
	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType(ext)

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("warning: failed to read config file: %s, will use environment variables", err.Error())
	}
}

// GetString retrieves a string config value from viper, falls back to os.Getenv if empty
func GetString(key string) string {
	value := viper.GetString(key)
	if value == "" {
		value = os.Getenv(key)
	}
	return value
}

// GetInt retrieves an int config value from viper, falls back to os.Getenv if viper returns 0
func GetInt(key string) int {
	value := viper.GetInt(key)
	if value == 0 {
		envValue := os.Getenv(key)
		if envValue != "" {
			intValue, err := strconv.Atoi(envValue)
			if err == nil {
				return intValue
			}
		}
	}
	return value
}

// GetBool retrieves a bool config value from viper, falls back to os.Getenv if not set in viper
func GetBool(key string) bool {
	if viper.IsSet(key) {
		return viper.GetBool(key)
	}
	envValue := os.Getenv(key)
	if envValue != "" {
		boolValue, err := strconv.ParseBool(envValue)
		if err == nil {
			return boolValue
		}
	}
	return false
}
