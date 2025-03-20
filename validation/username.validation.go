package validation

import (
	"regexp"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

func ValidateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// Check length (3 to 20 characters)
	if utf8.RuneCountInString(username) < 3 || utf8.RuneCountInString(username) > 20 {
		return false
	}

	// Check for allowed characters (alphanumeric, underscore, hyphen)
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validUsername.MatchString(username) {
		return false
	}

	// Check if username starts or ends with an underscore or hyphen
	if username[0] == '_' || username[0] == '-' || username[len(username)-1] == '_' || username[len(username)-1] == '-' {
		return false
	}

	return true
}
