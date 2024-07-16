package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost      string `mapstructure:"POSTGRES_HOST"`
	DBUsername  string `mapstructure:"POSTGRES_USER"`
	DBPassword  string `mapstructure:"POSTGRES_PASSWORD"`
	DBName      string `mapstructure:"POSTGRES_DB"`
	DBPort      string `mapstructure:"POSTGRES_PORT"`
	KongUrl     string `mapstructure:"KONG_URL"`
	SwaggerHost string `mapstructure:"SWAGGER_HOST"`
	SwaggerUrl  string `mapstructure:"SWAGGER_URL"`
	Environment string `mapstructure:"ENVIRONMENT"`
	ServerPort  string `mapstructure:"PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)

	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
