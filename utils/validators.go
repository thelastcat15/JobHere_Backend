package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("e164", validateE164)
	validate.RegisterValidation("latitude", validateLatitude)
	validate.RegisterValidation("longitude", validateLongitude)
}

func ValidateStruct(data interface{}) []ValidationErrorMsg {
	var errors []ValidationErrorMsg

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationErrorMsg{
				Field:   err.Field(),
				Message: err.Tag(),
			})
		}
	}

	return errors
}

func validateE164(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true
	}
	e164Regex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return e164Regex.MatchString(phone)
}

func validateLatitude(fl validator.FieldLevel) bool {
	lat := fl.Field().String()
	latRegex := regexp.MustCompile(`^-?([0-8]?[0-9]|90)(\.[0-9]{1,8})?$`)
	return latRegex.MatchString(lat)
}

func validateLongitude(fl validator.FieldLevel) bool {
	lon := fl.Field().String()
	lonRegex := regexp.MustCompile(`^-?([0-9]|[1-9][0-9]|1[0-7][0-9]|180)(\.[0-9]{1,8})?$`)
	return lonRegex.MatchString(lon)
}
