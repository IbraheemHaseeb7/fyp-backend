package validation

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func IsComsatsRegistrationNumber(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}

	pattern := `^(FA|SP)\d{2}-[a-zA-Z]{3}-\d{3}$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(fl.Field().String())
}

func IsPassword(fl validator.FieldLevel) bool {

	s := fl.Field().String()
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
