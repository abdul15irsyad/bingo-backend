package config

import (
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

var (
	Port         int
	CookieDomain string

	JWTPrivateKey string
	JWTPublicKey  string

	DBConfig DatabaseConfig
)

func InitAppConfig() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4020"
	}
	Port, _ = strconv.Atoi(port)

	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	CookieDomain = cookieDomain
	if CookieDomain == "" {
		CookieDomain = "localhost"
	}

	privateKey := os.Getenv("JWT_PRIVATE_KEY")
	publicKey := os.Getenv("JWT_PUBLIC_KEY")
	JWTPrivateKey = privateKey
	JWTPublicKey = publicKey

	DBConfig = DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
	}

	return nil
}
