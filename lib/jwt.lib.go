package lib

import (
	"bingo/config"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(sub string) (string, error) {
	privateKey, _ := parsePrivateKey(config.JWTPrivateKey)
	claims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(5 * time.Hour)},
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseJWT(tokenString string) (sub string, err error) {
	publicKey, _ := parsePublicKey(config.JWTPublicKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return publicKey, nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("failed to parse claims")
	}

	sub, err = claims.GetSubject()
	if err != nil {
		return "", err
	}

	return sub, nil
}

func parsePrivateKey(privateKeyString string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(strings.ReplaceAll(privateKeyString, `\n`, "\n")))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key format")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func parsePublicKey(publicKeyString string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(strings.ReplaceAll(publicKeyString, `\n`, "\n")))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid public key format")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, err
	}

	return rsaPublicKey, nil
}
