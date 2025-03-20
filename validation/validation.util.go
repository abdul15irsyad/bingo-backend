package validation

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   any    `json:"value"`
}

var (
	Validator *validator.Validate
	Trans     ut.Translator
)

func InitValidation() {
	Validator = validator.New(validator.WithRequiredStructEnabled())
	Validator.RegisterValidation("username", ValidateUsername)

	idn := id.New()
	uni := ut.New(idn, idn)
	trans, _ := uni.GetTranslator("id")
	Trans = trans
	id_translations.RegisterDefaultTranslations(Validator, Trans)

	Validator.RegisterTranslation(
		"username",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("username", "{0} must be 3-20 characters, contain only letters, numbers, underscores, or hyphens, and cannot start or end with an underscore or hyphen", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("username", fe.Tag())
			return t
		})
}

func Validate[T any](c *gin.Context, dtos T) []ValidationError {
	validationErrors := []ValidationError{}
	err := Validator.Struct(dtos)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.StructField()
			field, _ := reflect.TypeOf(dtos).FieldByName(fieldName)
			tagValue := field.Tag.Get("json")
			if tagValue == "" {
				tagValue = err.Field()
			}
			validationError := ValidationError{
				Field:   tagValue,
				Tag:     err.Tag(),
				Value:   err.Value(),
				Message: err.Translate(Trans),
			}
			validationErrors = append(validationErrors, validationError)
		}
	}
	return validationErrors
}
