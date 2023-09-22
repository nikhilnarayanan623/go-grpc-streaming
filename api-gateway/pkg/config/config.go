package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	ApiPort           string `mapstructure:"API_PORT"`
	StreamServiceHost string `mapstructure:"STREAMER_SERVICE_HOST"`
	StreamServicePort string `mapstructure:"STREAMER_SERVICE_PORT"`
}

var envs = []string{"API_PORT", "STREAMER_SERVICE_HOST", "STREAMER_SERVICE_PORT"}

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
