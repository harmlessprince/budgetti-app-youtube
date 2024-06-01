package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/harmlessprince/bougette-backend/common"
	"github.com/labstack/echo/v4"
	"reflect"
	"strings"
)

func (h *Handler) ValidateBodyRequest(c echo.Context, payload interface{}) []*common.ValidationError {
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())
	var errors []*common.ValidationError
	err := validate.Struct(payload)
	validationErrors, ok := err.(validator.ValidationErrors)
	if ok {
		reflected := reflect.ValueOf(payload)

		for _, validationErr := range validationErrors {
			field, _ := reflected.Type().FieldByName(validationErr.StructField())

			key := field.Tag.Get("json")
			if key == "" {
				key = strings.ToLower(validationErr.StructField())
			}
			condition := validationErr.Tag()
			keyToTitleCase := strings.Replace(key, "_", " ", -1)
			errMessage := keyToTitleCase + " field is " + condition

			switch condition {
			case "required":
				errMessage = keyToTitleCase + " is required"
			case "email":
				errMessage = keyToTitleCase + " must be a valid email address"
			}

			currentValidationError := &common.ValidationError{
				Error:     errMessage,
				Key:       key,
				Condition: condition,
			}
			errors = append(errors, currentValidationError)
		}
	}

	return errors
}
