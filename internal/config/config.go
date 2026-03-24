package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL          string
	MidtransServerKey    string
	MidtransClientKey    string
	MidtransIsProduction bool
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No .env file found")
	}

	config := &Config{
		DatabaseURL:          viper.GetString("DATABASE_URL"),
		MidtransServerKey:    viper.GetString("MIDTRANS_SERVER_KEY"),
		MidtransClientKey:    viper.GetString("MIDTRANS_CLIENT_KEY"),
		MidtransIsProduction: viper.GetBool("MIDTRANS_IS_PRODUCTION"),
	}
	return config, nil
}
