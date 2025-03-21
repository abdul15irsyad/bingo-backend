package config

import (
	"os"
	"strconv"
)

var (
	Port         int
	CookieDomain string

	JWTPrivateKey string
	JWTPublicKey  string
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

	return nil
}
