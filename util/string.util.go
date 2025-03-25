package util

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strings"
)

func Slugify(text string) string {
	slug := strings.ToLower(text)

	re := regexp.MustCompile(`[^\w\s-]`)
	slug = re.ReplaceAllString(slug, "")

	re = regexp.MustCompile(`[\s-]+`)
	slug = re.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}

var defaultLetters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int, letters *string) (string, error) {
	if letters == nil {
		letters = &defaultLetters
	}
	ret := make([]byte, n)
	for i := range n {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(*letters))))
		if err != nil {
			return "", err
		}
		ret[i] = (*letters)[num.Int64()]
	}
	return string(ret), nil
}
