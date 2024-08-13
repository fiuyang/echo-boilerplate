package config

import (
	"github.com/joho/godotenv"
	"os"
)

func Get() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	return &Config{
		Server: Server{
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: Database{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Name: os.Getenv("DB_NAME"),
		},
		Swagger: Swagger{
			Host: os.Getenv("SWAGGER_HOST"),
			Url:  os.Getenv("SWAGGER_URL"),
			Mode: os.Getenv("SWAGGER_MODE"),
		},
		Kong: Kong{
			Url: os.Getenv("KONG_URL"),
		},
		Obs: ObsHuawei{
			Ak:       os.Getenv("OBS_HUAWEI_AK"),
			Sk:       os.Getenv("OBS_HUAWEI_SK"),
			Endpoint: os.Getenv("OBS_HUAWEI_ENDPOINT"),
			Bucket:   os.Getenv("OBS_HUAWEI_BUCKET"),
		},
	}
}
