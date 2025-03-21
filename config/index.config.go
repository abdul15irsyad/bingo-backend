package config

import (
	"github.com/joho/godotenv"
)

func InitConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	InitAppConfig()

	return nil
}
