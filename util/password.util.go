package util

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	cost := bcrypt.DefaultCost // Adjust cost as needed (higher = more secure, slower)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hashedPassword), err
}

func ComparePassword(hashedPassword string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
