package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	StreamServiceHost string `mapstructure:"STREAMER_SERVICE_HOST"`
	StreamServicePort string `mapstructure:"STREAMER_SERVICE_PORT"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBName            string `mapstructure:"DB_NAME"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
}

var envs = []string{
	"STREAMER_SERVICE_HOST", "STREAMER_SERVICE_PORT",
	"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD",
}

func LoadConfig() (Config, error) {
	var config Config

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}

	return config, nil
}
