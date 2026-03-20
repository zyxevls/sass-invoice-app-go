package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No .env file found")
	}

	config := &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
	}
	return config, nil
}
